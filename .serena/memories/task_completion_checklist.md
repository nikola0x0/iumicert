# Task Completion Checklist

## When Implementing New Features

### Backend (Go) Changes
1. **Code Quality**
   - Run `go fmt` on modified files
   - Ensure error handling follows project patterns
   - Add appropriate logging for debugging
   - Update CLI help text if adding commands

2. **Testing**
   - Test CLI commands manually: `./micert <command> --help`
   - Verify API endpoints work: `curl http://localhost:8080/api/<endpoint>`
   - Check system health: `curl http://localhost:8080/api/health`

3. **Data Pipeline**
   - If modifying data flow, run full pipeline: `./reset.sh && ./generate.sh`
   - Verify receipts generate correctly
   - Test local verification: `./micert verify-local <receipt>`

### Frontend (React/Next.js) Changes
1. **Development**
   - Run `npm run lint` to check code style
   - Ensure TypeScript compilation: `npm run build`
   - Test in development mode: `npm run dev`

2. **Integration**
   - Verify API integration works with backend
   - Test Web3 functionality if wallet features modified
   - Check responsive design on different screen sizes

### Blockchain Integration
1. **Environment Setup**
   - Ensure `.env` has `ISSUER_PRIVATE_KEY` set
   - Test Sepolia connectivity
   - Verify gas estimation works

2. **Smart Contract Interaction**
   - Test publishing: `./micert publish-roots`
   - Verify transaction appears on Sepolia etherscan
   - Check contract interaction in frontend

## Final Verification Steps

### System Health Check
```bash
# Backend health
curl http://localhost:8080/api/health

# CLI functionality  
./micert version
./micert --help

# Data integrity
ls data/student_journeys/students/
ls publish_ready/receipts/
```

### Integration Testing
```bash
# Full pipeline test
./reset.sh && ./generate.sh

# API + Frontend
cd web/iumicert-issuer && npm run dev
# Open http://localhost:3000 and test UI

# Blockchain (if applicable)
./micert publish-roots --help
```

### Documentation Updates
- Update README.md if adding new features
- Update CLI help text for new commands
- Add API endpoint documentation if needed
- Update environment variable documentation