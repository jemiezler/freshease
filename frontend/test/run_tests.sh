#!/bin/bash
# test/run_tests.sh

echo "🧪 Running FreshEase Frontend Tests"
echo "=================================="

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Change to project directory
cd "$PROJECT_DIR"

echo "📦 Installing dependencies..."
flutter pub get

echo "🔧 Generating mock files..."
flutter packages pub run build_runner build --delete-conflicting-outputs

echo "🧪 Running unit tests..."
flutter test test/unit/

echo "🎨 Running widget tests..."
flutter test test/widgets/

echo "🔗 Running integration tests..."
flutter test integration_test/

echo "📊 Running all tests with coverage..."
flutter test --coverage

echo "✅ All tests completed!"
echo "📈 Coverage report generated in coverage/lcov.info"
