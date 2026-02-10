#!/bin/bash

echo "🚀 Setting up NowYouSeeMe Demo..."
echo ""

# Check if backend is running
echo "📡 Checking backend server..."
if curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "✓ Backend is running"
else
    echo "✗ Backend is not running!"
    echo ""
    echo "Please start the backend first:"
    echo "  cd backend"
    echo "  go run main.go"
    echo ""
    exit 1
fi

echo ""
echo "🎨 Generating sample visualizations..."
cd sdk
python3 examples/generate_sample_data.py

echo ""
echo "✨ Demo setup complete!"
echo ""
echo "🌐 Now start the frontend:"
echo "  cd frontend"
echo "  npm run dev"
echo ""
echo "Then visit http://localhost:3000"
