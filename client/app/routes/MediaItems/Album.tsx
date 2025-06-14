import { useState } from "react";
import { useLoaderData, type LoaderFunctionArgs } from "react-router";
import TopTracks from "~/components/TopTracks";
import { mergeAlbums, type Album } from "api/api";
import LastPlays from "~/components/LastPlays";
import PeriodSelector from "~/components/PeriodSelector";
import MediaLayout from "./MediaLayout";
import ActivityGrid from "~/components/ActivityGrid";

export async function clientLoader({ params }: LoaderFunctionArgs) {
  const res = await fetch(`/apis/web/v1/album?id=${params.id}`);
  if (!res.ok) {
    throw new Response("Failed to load album", { status: 500 });
  }
  const album: Album = await res.json();
  return album;
}

export default function Album() {
  const album = useLoaderData() as Album;
  const [period, setPeriod] = useState('week')

  console.log(album)

  return (
    <MediaLayout type="Album"
        title={album.title}
        img={album.image}
        id={album.id}
        musicbrainzId={album.musicbrainz_id}
        imgItemId={album.id}
        mergeFunc={mergeAlbums}
        mergeCleanerFunc={(r, id) => {
            r.artists = []
            r.tracks = []
            for (let i = 0; i < r.albums.length; i ++) {
                if (r.albums[i].id === id) {
                    delete r.albums[i]
                }
            }
            return r
        }}
        subContent={<>
        {album.listen_count && <p>{album.listen_count} play{ album.listen_count > 1 ? 's' : ''}</p>}
        </>}
    >
        <div className="mt-10">
            <PeriodSelector setter={setPeriod} current={period} />
        </div>
        <div className="flex flex-wrap gap-20 mt-10">
            <LastPlays limit={30} albumId={album.id} />
            <TopTracks limit={12} period={period} albumId={album.id} />
            <ActivityGrid autoAdjust configurable albumId={album.id} />
        </div>
    </MediaLayout>
  );
}
