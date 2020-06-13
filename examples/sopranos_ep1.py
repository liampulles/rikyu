from pathlib import Path

from rikyu import Project, Seq, Uniq

proj = Project(Path("/media/liam/Additional 2/Encodes/testdir/"))

# TODO: Rerun for each disc, then comment out
proj.rip_from_folder(Path("/media/liam/SOPRANOS_SEASON1_DISC1/"))\
    .title(1, Seq("episode", 1))\
    .title(3, Seq("episode", 2))\
    .title(4, Seq("episode", 3))\
    .title(5, Uniq("making of"))\
    .rip()

# TODO: Work on episodes
# ep_pipe = pipeline([
#     step("human_subtitle_complete", human(providing="subtitles")),
#     step("extract", vob_extract(audio=True, chapters=True)),
#     step("human_chapter_complete", human(providing="chapters"), depends_on=["extract"]),
#     step("video", x264(crf=21.0, preset="very_slow", tune="film", vapoursynth=vap), depends_on=["extract"]),
#     step("filter_audio", filter_(has_language("en") and is_largest()) or has_content("commentary"))),
#     step("audio", eac3to(codec="opus", slowdown=True, dpl2=True), depends_on=["extract"]),
#     step("fetch_metadata", imdb(series="Sopranos", season=1)),
#     step("mux", mkvmerge(video_fps=23.976), depends_on=["human_subtitle_complete", "human_chapter_complete", "video", "audio", "fetch_metadata"])
#     ])
# proj.apply_pipeline(ep_pipe, [seq("episode", i) for i in range(1,14)])

# ...and similar for extras.

# TODO: --- RIP ---
    # TODO: Rip all the dvd contents to their folders

# TODO: --- FIRST EPISODE ---
    # TODO: Work on the vapoursynth
    # TODO: Setup pipeline for audio video and container.
    # TODO: Do a "compressibility check", tune x264 and vapoursynth accordingly
    # TODO: When happy, run for full episode.
    # TODO: Review, tune and rerun if necessary until happy.

# TODO: -- REST OF EPISODES ---
    # TODO: Create batch pipeline for all episodes
    # TODO: As episodes complete, if they need individual adjustments - move to their own sections with special settings.
    # TODO: While they're all going, I'm going to work on the subtitles (there should be a special "human interaction"
    #       section created upfront from the above for me to use)
    # TODO: Similar for the chapters, though it should give me the base text file (might be able to automate, but fine for now).

# TODO: --- EXTRAS ---
    # TODO: Pretty much the same as for episodes, though less likely I'll need to redo them and I probably won't
    #       do chapters.
