FROM mcr.microsoft.com/dotnet/core/sdk:3.1 AS build-env
RUN git clone https://github.com/obecto/perper /perper
COPY src /apocryph/src
RUN cd /apocryph/src/Apocryph.Runtime.FunctionApp && \
    mkdir -p /home/site/wwwroot && \
    dotnet publish *.csproj --output /home/site/wwwroot

FROM mcr.microsoft.com/azure-functions/dotnet:3.0
ENV AzureWebJobsScriptRoot=/home/site/wwwroot \
    AzureFunctionsJobHost__Logging__Console__IsEnabled=true \
    ASPNETCORE_URLS=http://*:7071
COPY --from=build-env ["/home/site/wwwroot", "/home/site/wwwroot"]