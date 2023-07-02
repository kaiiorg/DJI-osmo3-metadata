# DJI Osmo Action 3 Metadata Reverse Engineering Notes
Got bored and wondered if my new DJI Osmo Action 3 recorded any extra datastreams to the files it saves. It turns out it does. These are some notes, code, and data dumps of a video recorded by said camera.

Long story short, the extra data channels appear to contain protocol buffer data, but I have no way of decoding them to something useful. Maybe if I knew of a program DJI released that actually uses this data, it could be further reverse engineered by comparing its outputs with the raw data. Or, even better, DJI could release the protobuf definition files and we can read them directly. That'll never happen.

Mind the ugly code.

## ffprobe
```
Input #0, mov,mp4,m4a,3gp,3g2,mj2, from 'DJI_20230701103032_0001_D.MP4':
  Metadata:
    major_brand     : isom
    minor_version   : 512
    compatible_brands: isomiso2avc1mp41
    creation_time   : 2023-07-01T15:30:37.000000Z
    encoder         : DJI OsmoAction3
    timecode        : 11:47:29;28
  Duration: 00:02:14.78, start: 0.000000, bitrate: 100393 kb/s
    Stream #0:0(und): Video: h264 (High) (avc1 / 0x31637661), yuv420p(tv, bt709), 2688x1512, 95663 kb/s, 59.94 fps, 59.94 tbr, 60k tbn, 119.88 tbc (default)
    Metadata:
      creation_time   : 2023-07-01T15:30:37.000000Z
      handler_name    : VideoHandler
      timecode        : 11:47:29;28
    Stream #0:1(und): Audio: aac (LC) (mp4a / 0x6134706D), 48000 Hz, stereo, fltp, 317 kb/s (default)
    Metadata:
      creation_time   : 2023-07-01T15:30:37.000000Z
      handler_name    : SoundHandler
    Stream #0:2(und): Data: none (djmd / 0x646D6A64), 26 kb/s
    Metadata:
      creation_time   : 2023-07-01T15:30:37.000000Z
      handler_name    : DJI meta
    Stream #0:3(und): Data: none (dbgi / 0x69676264), 4322 kb/s
    Metadata:
      creation_time   : 2023-07-01T15:30:37.000000Z
      handler_name    : DJI dbgi
    Stream #0:4(und): Data: none (tmcd / 0x64636D74)
    Metadata:
      creation_time   : 2023-07-01T15:30:37.000000Z
      handler_name    : TimeCodeHandler
      timecode        : 11:47:29;28
    Stream #0:5: Video: mjpeg (Baseline), yuvj420p(pc, bt470bg/unknown/unknown), 1280x720 [SAR 1:1 DAR 16:9], 90k tbr, 90k tbn, 90k tbc (attached pic)
```

- Contents of stream `0:2` dumped to `dji_meta.bin`
- Contents of stream `0:3` dumped to `dji_dbgi.bin`

## Notes
- Viewing `dji_meta.bin` and `dji_dbgi.bin` in a plain text editor shows mostly binary data, the first human readable text hints at usage of protocol buffers:
    - string `dvtm_ac002.proto` in `dji_meta.bin`
    - string `dbginfo_ac202.proto` in `dji_dbgi.bin`
- Unmarshalling a protobuf without the original definition files is really hard, but you can kinda do it
    - Doing that in main.go
- Confirmed files are protobuf data!
- Determined the top level data structure:
    - `dvtm_ac002.proto`
        - description; found single field of this type; see `dji_meta.1.json`
            - source protobuf definition
                - "dvtm_ac002.proto"
            - version number 1
                - "02.00.01"
                - Not `Firmware Version Number` as defined in `Device Info` menu (was "01.03.10.30")
                - Not `Camera Firmware Version Number` as defined in `Device Info` menu (was "10.00.31.12")
            - version number 2
                - "2.0.0"
                - Not `Firmware Version Number` as defined in `Device Info` menu (was "01.03.10.30")
                - Not `Camera Firmware Version Number` as defined in `Device Info` menu (was "10.00.31.12")
            - int value
                - 7500850
            - camera make/model
                - "DJI OsmoAction3"
        - video info; found single field of this type; see `dji_meta.2.json`
            - description string
                - "video"
            - video horizontal resolution
                - 2688
                - matches source video file horizontal resolution
            - video vertical resolution
                - 1512
                - matches source video file vertical resolution
            - int32 value 1
                - 1114620559
            - int value 1
                - 8
            - int value 2
                - 4
        - array of sample data; found 8077 fields of this type; see `dji_meta.3.json`
            - No idea what any of these values mean, not going to list them here
    - `dbginfo_ac202.proto`
        - description; found single field of this type; see `dji_dbgi.1.json`
            - int 1
                - 7500850
            - source protobuf definition
                - "dbginfo_ac202.proto"
            - version number 1
                - "3.1.0"
                - Not `Firmware Version Number` as defined in `Device Info` menu (was "01.03.10.30")
                - Not `Camera Firmware Version Number` as defined in `Device Info` menu (was "10.00.31.12")
                - Maybe version of protoc used by DJI?
            - version number 2
                - "2.0.2"
                - Not `Firmware Version Number` as defined in `Device Info` menu (was "01.03.10.30")
                - Not `Camera Firmware Version Number` as defined in `Device Info` menu (was "10.00.31.12")
                - Maybe version of protoc plugin used to compile the protobuf definition files to the source language used by DJI
            - sensor name
                - "IMX686"
                - Google results suggest its a Sony sensor
        - array of sample data; found 8077 instances of this field; see `dji_dbgi.2.json`
            - No idea what any of these values mean, not going to list them here