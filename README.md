## Test task - a wrapper for a Telegram bot.

Requirements:

+ Docker

This service allows you to send and receive text and file messages in Telegram. The service stores all messages in a
__PostgreSQL__ that will be run inside the docker automatically. In this way, you can get the entire message history.
If desired, you can specify a folder for storing data outside the container. https://prnt.sc/ptj2uCqSK0Tc

The service is built in such a way that you only need to enter a token and run one command. Docker will download all
dependencies by itself.

For ease of API testing, I have added swagger where you can try and see all API requests without any helper programs.

Sequence of actions:

1. Paste your token in the Dockerfile (in the BOT_TOKEN parameter).
2. Run the docker compose up --build command
3. After successful launch, go to swagger: http://localhost/swagger/index.html
4. Use necessary methods 🙂

Used libraries:
 + GORM
 + Goose (for creating enum type)
 + Swagger
 + GIN
 + Telegram