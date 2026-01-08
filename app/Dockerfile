# use the official Bun image
# see all versions at https://hub.docker.com/r/oven/bun/tags
FROM oven/bun:1 AS staging
WORKDIR /app

COPY package.json bun.lock* ./

# use --ignore-scripts to avoid building node modules like better-sqlite3
RUN bun install



# copy dependencies and source code into final image
FROM staging AS development

# RUN apt-get update && apt-get install -y python3 build-essential && ln -sf python3 /usr/bin/python
# RUN npm config set python "/usr/bin/python"

# run the app
COPY . .
EXPOSE 3000/tcp
ENTRYPOINT [ "bun", "--bun", "run", "dev", "-o" ]


FROM staging AS build
# Copy the entire project
COPY . .

RUN bun --bun run build

# copy production dependencies and source code into final image
FROM oven/bun:1 AS production
WORKDIR /app

# Only `.output` folder is needed from the build stage
COPY --from=build /app/.output /app

# run the app
EXPOSE 3000/tcp
ENTRYPOINT [ "bun", "--bun", "run", "/app/server/index.mjs" ]
