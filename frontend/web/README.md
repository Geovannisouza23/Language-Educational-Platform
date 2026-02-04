# Language Learning Platform - Web Frontend

Next.js 14 web application for the Language Learning Platform.

## Features

- 🎨 Modern UI with Tailwind CSS
- ⚡ Next.js 14 with App Router
- 🔐 JWT Authentication
- 📱 Responsive Design
- 🌐 API Integration
- 🎯 TypeScript Support

## Getting Started

### Prerequisites
- Node.js 20+
- npm or yarn

### Installation

```bash
cd frontend/web
npm install
```

### Environment Variables

Create a `.env.local` file:

```env
NEXT_PUBLIC_API_URL=http://localhost
```

### Development

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000)

### Build

```bash
npm run build
npm start
```

## Project Structure

```
web/
├── app/                 # App Router pages
│   ├── layout.tsx      # Root layout
│   ├── page.tsx        # Home page
│   └── globals.css     # Global styles
├── components/         # React components
├── lib/               # Utilities and services
│   ├── api.ts         # API client
│   └── courseService.ts
├── public/            # Static files
└── package.json
```

## Pages

- `/` - Home page
- `/courses` - Browse courses
- `/login` - Authentication
- `/dashboard` - User dashboard
- `/course/[id]` - Course details

## Technologies

- **Framework**: Next.js 14
- **Styling**: Tailwind CSS
- **State**: Zustand
- **Data Fetching**: React Query
- **HTTP Client**: Axios
- **Language**: TypeScript

## API Integration

All API calls go through the centralized API client in `lib/api.ts` with automatic token management.

## Deployment

Build the production bundle:

```bash
npm run build
```

Deploy to Vercel, Netlify, or any Node.js hosting platform.
