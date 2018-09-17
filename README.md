# youtubets

Downloading youtube transcript tool

*NOTE:* `youtubets` does not support automatic translation transcript.

## Usage

get transcript
```bash
youtubets [video_id]
```

`video_id` can get in Youtube URL.

ex.) https://www.youtube.com/watch?v=abcde => `video_id = abcde`

get transcript list

```bash
youtubets -l [video_id]
```

Youtube videos may have some transcripts for different languages.
List option show the list of transscripts. 

After showed the list, you can get specific language.

```bash
youtubets -lang en -name "English" [video_id]
```

## Installation
From the root of your cloned peco repository, run:
```bash
make deps
make install
```
