/*
 * Copyright 2026, Google LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package anthos.samples.bankofanthos.ledgerwriter;

import io.opentelemetry.api.common.AttributeKey;
import io.opentelemetry.api.common.Attributes;
import io.opentelemetry.exporter.otlp.http.logs.OtlpHttpLogRecordExporter;
import io.opentelemetry.instrumentation.log4j.appender.v2_17.OpenTelemetryAppender;
import io.opentelemetry.sdk.OpenTelemetrySdk;
import io.opentelemetry.sdk.logs.SdkLoggerProvider;
import io.opentelemetry.sdk.logs.export.BatchLogRecordProcessor;
import io.opentelemetry.sdk.resources.Resource;
import javax.annotation.PreDestroy;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.context.annotation.Configuration;

/**
 * Sends Log4j2 application logs to the Collector when {@code OTEL_EXPORTER_OTLP_ENDPOINT}
 * is set. Stdout logging stays in place. {@code gke-dev} leaves the variable unset.
 */
@Configuration
public class OtlpLogConfig {

    private static final Logger LOGGER = LogManager.getLogger(OtlpLogConfig.class);
    private static final String SERVICE_NAME = "ledgerwriter";

    private final SdkLoggerProvider loggerProvider;

    OtlpLogConfig() {
        this.loggerProvider = installIfEndpointSet();
    }

    @PreDestroy
    void shutdown() {
        if (loggerProvider != null) {
            loggerProvider.close();
        }
    }

    private static SdkLoggerProvider installIfEndpointSet() {
        String endpoint = System.getenv("OTEL_EXPORTER_OTLP_ENDPOINT");
        if (endpoint == null || endpoint.isBlank()) {
            return null;
        }
        String serviceName = System.getenv("OTEL_SERVICE_NAME");
        if (serviceName == null || serviceName.isBlank()) {
            serviceName = SERVICE_NAME;
        }
        Resource resource = Resource.getDefault().merge(Resource.create(
                Attributes.of(AttributeKey.stringKey("service.name"), serviceName)));
        SdkLoggerProvider provider = SdkLoggerProvider.builder()
                .setResource(resource)
                .addLogRecordProcessor(BatchLogRecordProcessor.builder(
                        OtlpHttpLogRecordExporter.builder().build()).build())
                .build();
        OpenTelemetryAppender.install(OpenTelemetrySdk.builder()
                .setLoggerProvider(provider)
                .build());
        LOGGER.info("OTLP logs enabled endpoint={} service={}", endpoint.trim(), serviceName);
        return provider;
    }
}
