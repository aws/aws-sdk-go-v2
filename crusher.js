const fs = require('fs');

const files = [
    "greengrass.json",
    "amplifybackend.json",
    "mediaconnect.json",
    "route53-recovery-control-config.json",
    "pinpoint.json",
    "apigatewayv2.json",
    "mediaconvert.json",
    "medialive.json",
    "macie2.json",
    "mediapackage.json",
    "apigatewaymanagementapi.json",
    "kafka.json",
    "mediapackage-vod.json",
    "mq.json",
    "iot-1click-devices-service.json",
    "serverlessapplicationrepository.json",
    "schemas.json",
    "pinpoint-sms-voice.json",
    "route53-recovery-readiness.json",
    "dataexchange.json",
    "kafkaconnect.json",
    "mediatailor.json"
];

const prefix = "codegen/sdk-codegen/aws-models";

const modelFiles = files.map(f => `${prefix}/${f}`);

for (const modelFile of modelFiles) {
    console.log(modelFile);
    const ast = JSON.parse(fs.readFileSync(modelFile).toString());
    fs.writeFileSync(modelFile, JSON.stringify(ast, hider, 4));
}

function hider(key, value) {
    if (key === "smithy.api#default" && value !== "") {
        return undefined;
    }
    return value;
}
