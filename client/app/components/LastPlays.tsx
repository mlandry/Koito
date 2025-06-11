import { useQuery } from "@tanstack/react-query"
import { timeSince } from "~/utils/utils"
import ArtistLinks from "./ArtistLinks"
import { getLastListens, type getItemsArgs } from "api/api"
import { Link } from "react-router"

interface Props {
    limit: number
    artistId?: Number
    albumId?: Number
    trackId?: number
    hideArtists?: boolean
}
  
export default function LastPlays(props: Props) {

    const { isPending, isError, data, error } = useQuery({ 
        queryKey: ['last-listens', {limit: props.limit, period: 'all_time', artist_id: props.artistId, album_id: props.albumId, track_id: props.trackId}], 
        queryFn: ({ queryKey }) => getLastListens(queryKey[1] as getItemsArgs),
    })

    if (isPending) {
        return (
            <div className="w-[500px]">
                <h2>Last Played</h2>
                <p>Loading...</p>
            </div>
        )
    }
    if (isError) {
        return <p className="error">Error:{error.message}</p>
    }

    let params = ''
    params += props.artistId ? `&artist_id=${props.artistId}` : ''
    params += props.albumId ? `&album_id=${props.albumId}` : ''
    params += props.trackId ? `&track_id=${props.trackId}` : ''

    return (
        <div>
            <h2 className="hover:underline"><Link to={`/listens?period=all_time${params}`}>Last Played</Link></h2>
            <table>
                <tbody>
                {data.items.map((item) => (
                    <tr key={`last_listen_${item.time}`}>
                        <td className="color-fg-tertiary pr-4 text-sm" title={new Date(item.time).toString()}>{timeSince(new Date(item.time))}</td>
                        <td className="text-ellipsis overflow-hidden max-w-[600px]">
                            {props.hideArtists ? <></> : <><ArtistLinks artists={item.track.artists} /> - </>}
                            <Link className="hover:text-(--color-fg-secondary)" to={`/track/${item.track.id}`}>{item.track.title}</Link>
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    )
}