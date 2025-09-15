# Servin Desktop GUI

A web interface for managing Servin containers, images, and volumes. Built with Flask backend and modern web frontend, packaged in a cross-platform desktop application.

## Features

- **Container Management**: View, start, stop, and remove containers
- **Image Management**: Browse, import, and manage container images  
- **Volume Management**: Create and manage persistent volumes
- **Real-time Monitoring**: Live updates of container status and logs
- **Cross-Platform**: Works on Windows, Linux, and macOS
- **Modern UI**: Dark theme with responsive design
- **Dual Mode**: Desktop app with embedded browser or standalone web interface

## Platform Support

| Platform | Architecture | Binary | Status |
|----------|-------------|---------|---------|
| Windows | amd64 | `servin.exe` | ✅ Supported |
| Linux | amd64 | `servin` | ✅ Supported |
| Linux | arm64 | `servin` | ✅ Supported |
| macOS | amd64 | `servin` | ✅ Supported |
| macOS | arm64 | `servin` | ✅ Supported |

## Quick Start

### Prerequisites

1. **Python 3.8+** with pip
2. **Servin binary** for your platform (auto-detected)

### Installation

1. **Clone and setup**:
   ```bash
   cd webview_gui
   python -m venv venv
   source venv/bin/activate  # Linux/macOS
   # or
   venv\Scripts\activate     # Windows
   pip install -r requirements.txt
   ```

2. **Build Servin for your platform**:
   ```bash
   # Windows
   .\build.ps1 -Target all
   
   # Linux/macOS  
   ./build-cross.sh --all
   ```

3. **Run the GUI**:
   ```bash
   # Desktop app with embedded browser
   python main.py
   
   # Web interface in default browser
   python demo.py
   ```

### Development Mode

For development and testing:

```bash
# Test all components
python test_app.py

# Run with mock data (no servin required)
python -c "import os; os.environ['SERVIN_MOCK'] = '1'; exec(open('demo.py').read())"
```

## Architecture

```
docker_gui/
├── app.py                 # Flask backend API
├── main.py               # Desktop app launcher  
├── demo.py               # Web demo launcher
├── servin_client.py      # Servin runtime interface
├── mock_servin_client.py # Mock client for demos
├── test_app.py          # Test suite
├── templates/
│   └── index.html       # Web interface
├── static/
│   ├── style.css        # Styling
│   └── script.js        # Frontend logic
└── requirements.txt     # Python dependencies
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Web interface |
| `/api/containers` | GET | List containers |
| `/api/containers/{id}/start` | POST | Start container |
| `/api/containers/{id}/stop` | POST | Stop container |
| `/api/containers/{id}/remove` | DELETE | Remove container |
| `/api/images` | GET | List images |
| `/api/images/{id}/remove` | DELETE | Remove image |
| `/api/volumes` | GET | List volumes |
| `/api/volumes` | POST | Create volume |
| `/api/system/info` | GET | System information |

## Configuration

The GUI automatically detects the servin binary using this priority:

1. **Platform-specific build directory** (`build/{platform}-{arch}/servin`)
2. **Root directory** (`servin` or `servin.exe`)
3. **Mock client** (for development/demo)

## Dependencies

- **Flask 3.0.3**: Web framework
- **Flask-CORS 4.0.1**: Cross-origin request handling
- **pywebview 5.1**: Embedded browser (optional)

## Development

### Adding Features

1. **Backend**: Add API endpoints in `app.py`
2. **Frontend**: Update `templates/index.html` and `static/script.js`
3. **Client**: Extend `servin_client.py` for new servin commands

### Testing

```bash
# Run full test suite
python test_app.py

# Test specific components
python -c "from app import app; print('Flask OK')"
python -c "from servin_client import ServinClient; print('Client OK')"
```

### Building for Distribution

```bash
# Build all platform binaries
.\build.ps1 -Target all        # Windows
./build-cross.sh --all         # Linux/macOS

# Package for distribution
# Creates: servin-{platform}-{version}.zip
```

## Troubleshooting

### Binary Not Found
- Ensure servin is built for your platform
- Check `build/{platform}-{arch}/` directory
- Run `python test_app.py` to verify detection

### Web Interface Issues
- Check Flask server is running on port 5555
- Verify firewall allows local connections
- Try demo mode: `python demo.py`

### pywebview Problems
- Install system dependencies for web engine
- Fallback to browser mode automatically
- Use `demo.py` for browser-only mode

## License

Same as Servin project license.

## Contributing

1. Fork the repository
2. Create feature branch
3. Add tests for new functionality
4. Submit pull request

For questions or issues, please open a GitHub issue.
