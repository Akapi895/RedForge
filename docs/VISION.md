# Vision Analysis (`analyze_image`)

## Overview

- **Tool name:** `analyze_image` (built-in MCP tool).
- **Behavior:** reads a local image, resizes and JPEG-compresses it with `imaging`, calls a separate **Vision** model, and returns **plain text** to the Agent.
- **Context:** image bytes are **not** written to conversation history; only the path and text summary enter the Agent context.


Vision analysis registers the `analyze_image` MCP tool when enabled. It is intended for screenshots, captchas, UI states, and image evidence in authorized workflows.

## Config

```yaml
vision:
  enabled: true
  model: qwen-vl
  api_key: ""
  base_url: ""
  provider: ""
  max_image_bytes: 5242880
  max_dimension: 2048
  jpeg_quality: 82
  max_payload_bytes: 524288
  skip_preprocess_below_bytes: 2097152  # pass through images below 2 MB when the long edge is within max_dimension; 0 always compresses to JPEG
  detail: low  # low | high | auto
  timeout_seconds: 60
```

Empty `api_key`, `base_url`, or `provider` inherits from the resolved default AI channel.
`model` is required. When `enabled: false`, the tool is not registered.

## Web Settings

**System Settings → Basic Settings → Vision Analysis (`analyze_image`)** configures the enable switch, vision model, API key/base URL (empty values reuse the default AI channel), and preprocessing parameters. **Save and apply** writes the values to `config.yaml` and registers the MCP tool again.

## Path Handling

`analyze_image` can read any readable image path on the server, either absolute or relative to the process working directory. The runtime still validates the image extension and regular-file type.

## Agent Usage

The system prompt instructs the Agent to call `analyze_image` for images and not use `read_file` to read binary image data.

`multi_agent.eino_middleware.tool_search_always_visible_tools` should include `analyze_image`.

**Data Handling**

Image bytes are sent only to the vision model call. Agent history keeps text summaries, not raw image bytes. This reduces context size and accidental image propagation.

**Preprocessing**

The runtime can resize and recompress large images based on:

- maximum file size;
- maximum dimension;
- JPEG quality;
- encoded payload size.

If small images are already under limits, preprocessing may be skipped.

**Usage Guidance**

Use vision for:

- UI screenshots;
- visual vulnerability evidence;
- captcha or image-based prompts in authorized tests;
- interpreting tool screenshots.

Do not use it for:

- unrelated personal images;
- sensitive screenshots without authorization;
- long-term storage of raw evidence when a text summary is enough.

## Compliance

When enabled, images are sent to the upstream Vision API configured for the tool. In sensitive environments, use a trusted gateway or keep `enabled: false`.

Source anchors: `internal/app/vision_tools.go`, `internal/vision/client.go`, `internal/vision/preprocess.go`, and `internal/config/vision.go`.
