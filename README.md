# Sambhav Copilot

The aim of this project is to build an assistant to help Interviewer conduct smooth interviews. Following are the features -
- Integrate with gmeet, and analyse audio transcription
- Create a detailed report in Notion docs
- Suggestion (Follow-up questions, deep dive questions, missed important resume-based questions, missed important job profile based questions) during interview using chrome extensions

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

Live reload the application:
```bash
make watch
```

Clean up binary from the last build:
```bash
make clean
```