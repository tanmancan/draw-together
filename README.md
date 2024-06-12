# draw-together

Demo: https://draw-together.dev/

## Quick-start using the command line tool

### Using the CLI

There is a CLI tool provided that can be used to build and run the application.

View all available options via `./draw-together -h`:

```
$ ./draw-together -h

Command line tool to help build and run the application.

Usage:

    ./draw-together [-abcdeprh]

Options:

    -b  Build the application, but don't run it.
        This is the default option.
    -c  Clear cache during build.
    -d  Build and run the application in development mode.
    -e  Setup .env file if they do not exists.
        Will create a .env and .env.dev for production and
        development builds.
    -p  Print the Bake options without building.
    -r  Build and run the application in production mode.
    -h  Show help text.
```

Generate `.env` file:

```bash
./draw-together -e
```

If you want to use the Azure image detection service replace the following values in `.env`:

```conf
# The Azure endpoint
AZURE_CV_ENDPOINT="https://example.cognitiveservices.azure.com/computervision"
# The API key
AZURE_CV_KEY="KEEP_IT_SECRET"
```

Build the app:

```bash
./draw-together -b
```

Run the app:

```bash
./draw-together -r
```

Application will be available via:

```
https://localhost:8443
```

Local development version uses a self signed certificate, so you may need to accept any browser warnings regarding SSL errors.
