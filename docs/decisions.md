# Technical Decisions

## ADR (Architecture Decision Records)

### ADR-001: Microservices Architecture

**Status:** Accepted

**Context:**
We need a scalable architecture that can handle multiple teachers, courses, and students concurrently.

**Decision:**
Implement microservices architecture with separate services for different domains (Auth, Users, Courses, Tasks, Progress).

**Consequences:**
- ✅ Independent scaling of services
- ✅ Technology diversity (C# for auth, Go for others)
- ✅ Isolated failures
- ⚠️ Increased complexity in deployment
- ⚠️ Distributed transaction challenges

---

### ADR-002: C# for Auth Service

**Status:** Accepted

**Context:**
Need robust authentication with proven libraries and strong type safety.

**Decision:**
Use C# with ASP.NET Core for Authentication Service.

**Rationale:**
- Mature JWT libraries
- Entity Framework for complex queries
- Strong type system
- Excellent IDE support
- Built-in security features

**Consequences:**
- ✅ Rapid development with proven patterns
- ✅ Strong community support
- ⚠️ Different stack from other services
- ⚠️ Requires .NET runtime

---

### ADR-003: Go for Business Services

**Status:** Accepted

**Context:**
Need high-performance, lightweight services for user, course, and task management.

**Decision:**
Use Go for all business logic services.

**Rationale:**
- Excellent performance
- Low memory footprint
- Fast compilation
- Built-in concurrency
- Single binary deployment
- Great for REST APIs

**Consequences:**
- ✅ High performance
- ✅ Easy deployment
- ✅ Consistent stack across services
- ⚠️ Less mature ORM ecosystem
- ⚠️ Simpler error handling

---

### ADR-004: PostgreSQL as Primary Database

**Status:** Accepted

**Context:**
Need reliable, ACID-compliant database with good performance.

**Decision:**
Use PostgreSQL for all services with separate databases per service.

**Rationale:**
- ACID compliance
- JSON support for flexibility
- Excellent performance
- Mature ecosystem
- Free and open-source

**Consequences:**
- ✅ Data consistency
- ✅ Rich feature set
- ✅ Good performance
- ⚠️ Separate database per service (data duplication)

---

### ADR-005: Redis for Caching

**Status:** Accepted

**Context:**
Need fast caching layer to reduce database load.

**Decision:**
Use Redis for caching user profiles, session data, and frequently accessed data.

**Rationale:**
- In-memory speed
- Rich data structures
- TTL support
- Pub/Sub capabilities
- Widely adopted

**Consequences:**
- ✅ Significantly improved performance
- ✅ Reduced database load
- ⚠️ Cache invalidation complexity
- ⚠️ Additional infrastructure component

---

### ADR-006: NGINX as API Gateway

**Status:** Accepted

**Context:**
Need single entry point for all services with rate limiting and SSL termination.

**Decision:**
Use NGINX as API Gateway.

**Rationale:**
- High performance
- SSL termination
- Rate limiting
- Load balancing
- Proven reliability

**Alternatives Considered:**
- Kong: Too complex for our needs
- Traefik: Good but NGINX more familiar
- AWS API Gateway: Vendor lock-in

**Consequences:**
- ✅ Single entry point
- ✅ Centralized security
- ✅ Rate limiting
- ⚠️ Single point of failure (mitigated with HA)

---

### ADR-007: JWT for Authentication

**Status:** Accepted

**Context:**
Need stateless authentication mechanism for distributed services.

**Decision:**
Use JWT (JSON Web Tokens) with refresh tokens.

**Rationale:**
- Stateless authentication
- Self-contained tokens
- Easy to validate across services
- Standard and widely supported

**Implementation:**
- Access token: 1 hour expiry
- Refresh token: 7 days expiry
- Tokens stored in Redis for revocation capability

**Consequences:**
- ✅ Scalable authentication
- ✅ No session storage needed
- ⚠️ Token size larger than session ID
- ⚠️ Cannot invalidate tokens immediately (mitigated with Redis blacklist)

---

### ADR-008: Docker and Kubernetes for Deployment

**Status:** Accepted

**Context:**
Need consistent deployment across environments with easy scaling.

**Decision:**
Use Docker for containerization and Kubernetes for orchestration.

**Rationale:**
- Consistent environments
- Easy scaling
- Self-healing
- Industry standard
- Great community support

**Consequences:**
- ✅ Consistent deployments
- ✅ Easy scaling
- ✅ Infrastructure as code
- ⚠️ Kubernetes complexity
- ⚠️ Requires DevOps expertise

---

### ADR-009: Monorepo Structure

**Status:** Accepted

**Context:**
Need to manage multiple services efficiently with shared code.

**Decision:**
Use monorepo structure with separate directories for each service.

**Rationale:**
- Single source of truth
- Shared libraries
- Coordinated changes
- Simplified CI/CD

**Consequences:**
- ✅ Easy to share code
- ✅ Atomic cross-service changes
- ✅ Simplified dependency management
- ⚠️ Larger repository size
- ⚠️ Need careful CI/CD configuration

---

### ADR-010: Observability with Prometheus and Grafana

**Status:** Accepted

**Context:**
Need comprehensive monitoring and alerting for production systems.

**Decision:**
Use Prometheus for metrics collection and Grafana for visualization.

**Rationale:**
- Industry standard
- Pull-based metrics
- Powerful query language (PromQL)
- Great visualization with Grafana
- Free and open-source

**Consequences:**
- ✅ Comprehensive monitoring
- ✅ Custom dashboards
- ✅ Alerting capabilities
- ⚠️ Additional infrastructure

---

### ADR-011: Database Per Service

**Status:** Accepted

**Context:**
Following microservices best practices, each service should own its data.

**Decision:**
Each service has its own database schema/instance.

**Rationale:**
- Service independence
- Independent scaling
- Technology flexibility
- Failure isolation

**Consequences:**
- ✅ Service autonomy
- ✅ Independent deployment
- ⚠️ Data duplication
- ⚠️ Distributed transactions complexity
- ⚠️ No foreign keys across services

---

## Future Decisions to Make

### Under Consideration:

1. **Message Queue Implementation**
   - RabbitMQ vs Kafka vs AWS SQS
   - For async communication between services

2. **Frontend Framework**
   - Next.js vs Remix
   - React Native vs Flutter for mobile

3. **Video Conferencing**
   - Zoom SDK vs Agora vs Jitsi
   - Self-hosted vs SaaS

4. **File Storage**
   - MinIO (self-hosted) vs AWS S3
   - Cost vs control trade-off

5. **Email Service**
   - SendGrid vs AWS SES vs Mailgun
   - Cost and deliverability considerations

6. **CI/CD Platform**
   - GitHub Actions (current) vs GitLab CI vs Jenkins
   - Cost and feature evaluation

7. **Cloud Provider**
   - AWS vs GCP vs Azure vs Self-hosted
   - Long-term cost analysis needed
