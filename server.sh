#!/bin/bash

# FlipAssistant Development Server Manager

case "$1" in
    "start"|"dev")
        echo "🚀 Starting FlipAssistant development servers..."
        make dev
        ;;
    "stop")
        echo "🛑 Stopping FlipAssistant servers..."
        make stop
        ;;
    "status")
        echo "📊 Checking server status..."
        make status
        ;;
    "restart")
        echo "🔄 Restarting FlipAssistant servers..."
        make stop
        sleep 2
        make dev
        ;;
    "install")
        echo "📦 Installing dependencies..."
        make install
        ;;
    "build")
        echo "🏗️  Building for production..."
        make build
        ;;
    "clean")
        echo "🧹 Cleaning up..."
        make clean
        ;;
    *)
        echo "FlipAssistant Server Manager"
        echo ""
        echo "Usage: ./server.sh [command]"
        echo ""
        echo "Commands:"
        echo "  start, dev  - Start development servers"
        echo "  stop        - Stop all servers"  
        echo "  status      - Check server status"
        echo "  restart     - Restart servers"
        echo "  install     - Install dependencies"
        echo "  build       - Build for production"
        echo "  clean       - Clean up and stop servers"
        echo ""
        echo "Examples:"
        echo "  ./server.sh start   # Start both servers"
        echo "  ./server.sh stop    # Stop all servers"
        echo "  ./server.sh status  # Check what's running"
        ;;
esac