# GenAI Health Integration

This document describes the integration of GenAI meal planning functionality with the health page in the FreshEase application.

## Overview

The integration connects real health data (steps, calories burned) from the device's health sensors with AI-powered meal plan generation. Users can generate personalized meal plans based on their daily activity levels.

## Architecture

### Backend (Go)
- **GenAI Module**: Located in `backend/modules/genai/`
- **API Endpoints**:
  - `POST /api/genai/daily` - Generate daily meal plan
  - `POST /api/genai/weekly` - Generate weekly meal plan
- **Integration**: Uses Google's Gemini AI to generate meal plans based on user profile and health data

### Frontend (Flutter)
- **GenAI Service**: `lib/core/genai/genai_service.dart`
- **API Client**: `lib/core/genai/genai_api.dart`
- **Data Models**: `lib/core/genai/models.dart`
- **UI Components**: `lib/core/genai/widgets.dart`
- **Health Controller**: Extended `lib/core/health/health_controller.dart`

## Features

### Real Data Integration
- **Steps**: Automatically fetched from device health sensors
- **Calories**: Calculated from active energy burned in the last 24 hours
- **User Profile**: Configurable gender, age, height, weight, and fitness goals

### AI Meal Planning
- **Personalized Plans**: Generated based on individual health metrics
- **Goal-Oriented**: Supports maintenance, weight loss, and weight gain targets
- **Real-time**: Uses current day's activity data for accurate recommendations

### User Interface
- **Health Dashboard**: Shows current steps and calories burned
- **Meal Plan Generator**: Interactive form for customization
- **Visual Meal Cards**: Clean display of generated meal plans with calorie breakdowns
- **Error Handling**: User-friendly error messages and loading states

## Usage

1. **Health Data Collection**: The app automatically requests health permissions and fetches daily activity data
2. **Meal Plan Generation**: Users can customize their profile and generate AI-powered meal plans
3. **Visualization**: Generated meal plans are displayed with detailed calorie information for each meal

## Data Flow

```
Device Health Sensors → Health Controller → GenAI Service → Backend API → Gemini AI → Meal Plan → UI Display
```

## Configuration

### Backend
- Requires `GOOGLE_API_KEY` environment variable for Gemini AI access
- GenAI module is registered as a public endpoint (no authentication required)

### Frontend
- GenAI service is injected via dependency injection
- Health controller manages both health data and meal plan generation
- Real-time updates through Flutter's ChangeNotifier pattern

## Benefits

- **Personalized**: Meal plans adapt to individual activity levels
- **Real-time**: Uses current day's health data for accurate recommendations
- **Integrated**: Seamless experience within the existing health dashboard
- **Visual**: Clear presentation of meal plans with nutritional information
