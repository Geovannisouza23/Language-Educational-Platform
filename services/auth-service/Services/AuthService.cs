using AuthService.Data;
using AuthService.Models;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Caching.Distributed;
using System.Security.Cryptography;

namespace AuthService.Services;

public class AuthService : IAuthService
{
    private readonly AuthDbContext _context;
    private readonly ITokenService _tokenService;
    private readonly IDistributedCache _cache;
    private readonly ILogger<AuthService> _logger;

    public AuthService(
        AuthDbContext context,
        ITokenService tokenService,
        IDistributedCache cache,
        ILogger<AuthService> logger)
    {
        _context = context;
        _tokenService = tokenService;
        _cache = cache;
        _logger = logger;
    }

    public async Task<AuthResponse> LoginAsync(LoginRequest request, string ipAddress, string userAgent)
    {
        var user = await _context.Users
            .Include(u => u.Role)
            .FirstOrDefaultAsync(u => u.Email == request.Email);

        if (user == null || !BCrypt.Net.BCrypt.Verify(request.Password, user.PasswordHash))
        {
            throw new UnauthorizedAccessException("Invalid credentials");
        }

        if (!user.IsActive)
        {
            throw new UnauthorizedAccessException("Account is inactive");
        }

        // Update last login
        user.LastLoginAt = DateTime.UtcNow;
        await _context.SaveChangesAsync();

        // Generate tokens
        var accessToken = _tokenService.GenerateAccessToken(user);
        var refreshToken = _tokenService.GenerateRefreshToken();
        
        var refreshTokenEntity = await _tokenService.SaveRefreshTokenAsync(
            user.Id, refreshToken, ipAddress, userAgent);

        _logger.LogInformation("User {Email} logged in successfully", user.Email);

        return new AuthResponse
        {
            AccessToken = accessToken,
            RefreshToken = refreshToken,
            ExpiresAt = refreshTokenEntity.ExpiresAt,
            User = new UserDto
            {
                Id = user.Id,
                Email = user.Email,
                Role = user.Role.Name,
                EmailVerified = user.EmailVerified
            }
        };
    }

    public async Task<AuthResponse> RegisterAsync(RegisterRequest request)
    {
        if (request.Password != request.ConfirmPassword)
        {
            throw new ArgumentException("Passwords do not match");
        }

        var existingUser = await _context.Users
            .FirstOrDefaultAsync(u => u.Email == request.Email);

        if (existingUser != null)
        {
            throw new InvalidOperationException("Email already registered");
        }

        var role = await _context.Roles
            .FirstOrDefaultAsync(r => r.Name == request.Role);

        if (role == null)
        {
            throw new ArgumentException("Invalid role");
        }

        var user = new User
        {
            Id = Guid.NewGuid(),
            Email = request.Email,
            PasswordHash = BCrypt.Net.BCrypt.HashPassword(request.Password),
            RoleId = role.Id,
            CreatedAt = DateTime.UtcNow,
            IsActive = true,
            EmailVerified = false
        };

        _context.Users.Add(user);
        await _context.SaveChangesAsync();

        _logger.LogInformation("New user registered: {Email}", user.Email);

        // Generate tokens
        var accessToken = _tokenService.GenerateAccessToken(user);
        var refreshToken = _tokenService.GenerateRefreshToken();
        
        var refreshTokenEntity = await _tokenService.SaveRefreshTokenAsync(
            user.Id, refreshToken, "127.0.0.1", "Registration");

        return new AuthResponse
        {
            AccessToken = accessToken,
            RefreshToken = refreshToken,
            ExpiresAt = refreshTokenEntity.ExpiresAt,
            User = new UserDto
            {
                Id = user.Id,
                Email = user.Email,
                Role = role.Name,
                EmailVerified = user.EmailVerified
            }
        };
    }

    public async Task<AuthResponse> RefreshTokenAsync(string refreshToken, string ipAddress, string userAgent)
    {
        var tokenEntity = await _context.RefreshTokens
            .Include(rt => rt.User)
            .ThenInclude(u => u.Role)
            .FirstOrDefaultAsync(rt => rt.Token == refreshToken);

        if (tokenEntity == null || tokenEntity.IsRevoked || tokenEntity.ExpiresAt < DateTime.UtcNow)
        {
            throw new UnauthorizedAccessException("Invalid or expired refresh token");
        }

        // Revoke old token
        tokenEntity.IsRevoked = true;
        
        // Generate new tokens
        var accessToken = _tokenService.GenerateAccessToken(tokenEntity.User);
        var newRefreshToken = _tokenService.GenerateRefreshToken();
        
        var newRefreshTokenEntity = await _tokenService.SaveRefreshTokenAsync(
            tokenEntity.UserId, newRefreshToken, ipAddress, userAgent);

        await _context.SaveChangesAsync();

        return new AuthResponse
        {
            AccessToken = accessToken,
            RefreshToken = newRefreshToken,
            ExpiresAt = newRefreshTokenEntity.ExpiresAt,
            User = new UserDto
            {
                Id = tokenEntity.User.Id,
                Email = tokenEntity.User.Email,
                Role = tokenEntity.User.Role.Name,
                EmailVerified = tokenEntity.User.EmailVerified
            }
        };
    }

    public async Task<bool> RevokeTokenAsync(string refreshToken)
    {
        var tokenEntity = await _context.RefreshTokens
            .FirstOrDefaultAsync(rt => rt.Token == refreshToken);

        if (tokenEntity == null)
        {
            return false;
        }

        tokenEntity.IsRevoked = true;
        await _context.SaveChangesAsync();

        return true;
    }

    public async Task<bool> VerifyEmailAsync(Guid userId, string token)
    {
        var user = await _context.Users.FindAsync(userId);
        if (user == null)
        {
            return false;
        }

        // Verify token from cache
        var cachedToken = await _cache.GetStringAsync($"email_verify_{userId}");
        if (cachedToken != token)
        {
            return false;
        }

        user.EmailVerified = true;
        user.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync();

        // Remove token from cache
        await _cache.RemoveAsync($"email_verify_{userId}");

        return true;
    }
}
