# Language Learning Platform - Mobile App

React Native mobile application for iOS and Android.

## Features

- 📱 Native iOS and Android apps
- 🎨 Modern UI components
- 🔐 JWT Authentication
- 🌐 API Integration
- 📲 Push Notifications (ready)
- 🎯 TypeScript Support

## Getting Started

### Prerequisites
- Node.js 18+
- React Native CLI
- Xcode (for iOS)
- Android Studio (for Android)

### Installation

```bash
cd frontend/mobile
npm install
```

### iOS Setup

```bash
cd ios
pod install
cd ..
npm run ios
```

### Android Setup

```bash
npm run android
```

## Project Structure

```
mobile/
├── App.tsx              # Main app component
├── src/
│   ├── api/            # API services
│   ├── components/     # React components
│   ├── screens/        # App screens
│   ├── navigation/     # Navigation setup
│   └── utils/          # Utilities
├── android/            # Android native code
├── ios/                # iOS native code
└── package.json
```

## Screens

- Home - Landing screen
- Courses - Browse courses
- CourseDetail - Course information
- MyLearning - Enrolled courses
- Profile - User profile
- VideoCall - Live classes

## API Integration

All API calls use the centralized client in `src/api/client.ts` with automatic token management via AsyncStorage.

## Build

### iOS
```bash
cd ios
xcodebuild -workspace LanguagePlatform.xcworkspace -scheme LanguagePlatform -configuration Release
```

### Android
```bash
cd android
./gradlew assembleRelease
```

## Technologies

- **Framework**: React Native 0.73
- **Navigation**: React Navigation 6
- **State**: Context API / Redux (optional)
- **Storage**: AsyncStorage
- **HTTP**: Axios
- **Language**: TypeScript

## Deployment

- **iOS**: Deploy via App Store Connect
- **Android**: Deploy via Google Play Console
