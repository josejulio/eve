EVE - Enhanced Virtual Engagement

1. Install go 1.22.3+
2. make
3. Update `config.yml` with your InstructLab's LLM service address
4. Set environment `OPENAI_API_KEY` to anything (just to bypass the lib checking)
5. Run `bin/eve`
6. Hit http: `http://localhost:8080/talk?input=how are you`
