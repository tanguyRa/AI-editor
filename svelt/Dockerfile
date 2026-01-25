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
EXPOSE 5173
ENTRYPOINT [ "bun", "run", "dev", "--host", "0.0.0.0" ]


FROM staging AS build
# Copy the entire project
COPY . .

RUN bun --bun run build
RUN bun --bun run prune --production

# copy production dependencies and source code into final image
FROM oven/bun:1 AS production
WORKDIR /app

# Only `.output` folder is needed from the build stage
COPY --from=build /app/build build/
COPY --from=build /app/node_modules node_modules/


# run the app
EXPOSE 5173
ENTRYPOINT [ "bun", "./build/index.js" ]
