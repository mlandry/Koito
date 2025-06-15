import { Link, useNavigate } from "react-router";
import ArtistLinks from "./ArtistLinks";
import { imageUrl, type Album, type Artist, type Track, type PaginatedResponse } from "api/api";

type Item = Album | Track | Artist;

interface Props<T extends Item> {
    data: PaginatedResponse<T>
    separators?: ConstrainBoolean
    type: "album" | "track" | "artist";
    className?: string,
}

export default function TopItemList<T extends Item>({ data, separators, type, className }: Props<T>) {

    return (
        <div className={`flex flex-col gap-1 ${className} min-w-[300px]`}>
            {data.items.map((item, index) => {
                const key = `${type}-${item.id}`;
                return (
                    <div
                        key={key}
                        style={{ fontSize: 12 }}
                        className={`${
                            separators && index !== data.items.length - 1 ? 'border-b border-(--color-fg-tertiary) mb-1 pb-2' : ''
                        }`}
                    >
                        <ItemCard item={item} type={type} key={type+item.id} />
                    </div>
                );
            })}
        </div>
    );
}

function ItemCard({ item, type }: { item: Item; type: "album" | "track" | "artist" }) {

    const itemClasses = `flex items-center gap-2 hover:text-(--color-fg-secondary)`

    const navigate = useNavigate();
    
    const handleItemClick = (type: string, id: number) => {
        navigate(`/${type.toLowerCase()}/${id}`);
    };
    
    const handleArtistClick = (event: React.MouseEvent) => {
        // Stop the click from navigating to the album page
        event.stopPropagation();
    };
    
    // Also stop keyboard events on the inner links from bubbling up
    const handleArtistKeyDown = (event: React.KeyboardEvent) => {
        event.stopPropagation();
    }

    switch (type) {
        case "album": {
            const album = item as Album;
    
            const handleKeyDown = (event: React.KeyboardEvent) => {
                if (event.key === 'Enter') {
                    handleItemClick("album", album.id);
                }
            };
    
            return (
                <div style={{fontSize: 12}}>
                    <div
                        className={itemClasses}
                        onClick={() => handleItemClick("album", album.id)}
                        onKeyDown={handleKeyDown}
                        role="link"
                        tabIndex={0}
                        aria-label={`View album: ${album.title}`}
                        style={{ cursor: 'pointer' }}
                    >
                        <img src={imageUrl(album.image, "small")} alt={album.title} />
                        <div>
                            <span style={{fontSize: 14}}>{album.title}</span>
                            <br />
                            {album.is_various_artists ?
                            <span className="color-fg-secondary">Various Artists</span>
                            :
                            <div onClick={handleArtistClick} onKeyDown={handleArtistKeyDown}>
                               <ArtistLinks artists={album.artists || [{id: 0, Name: 'Unknown Artist'}]}/>
                            </div>
                            }
                            <div className="color-fg-secondary">{album.listen_count} plays</div>
                        </div>
                    </div>
                </div>
            );
        }
        case "track": {
            const track = item as Track;
    
            const handleKeyDown = (event: React.KeyboardEvent) => {
                if (event.key === 'Enter') {
                    handleItemClick("track", track.id);
                }
            };

            return (
                <div style={{fontSize: 12}}>
                <div
                    className={itemClasses}
                    onClick={() => handleItemClick("track", track.id)}
                    onKeyDown={handleKeyDown}
                    role="link"
                    tabIndex={0}
                    aria-label={`View track: ${track.title}`}
                    style={{ cursor: 'pointer' }}
                >
                    <img src={imageUrl(track.image, "small")} alt={track.title} />
                    <div>
                        <span style={{fontSize: 14}}>{track.title}</span>
                        <br />
                            <div onClick={handleArtistClick} onKeyDown={handleArtistKeyDown}>
                               <ArtistLinks artists={track.artists || [{id: 0, Name: 'Unknown Artist'}]}/>
                            </div>
                        <div className="color-fg-secondary">{track.listen_count} plays</div>
                    </div>
                    </div>
                </div>
            );
        }
        case "artist": {
            const artist = item as Artist;
            return (
                <div style={{fontSize: 12}}>
                    <Link className={itemClasses+' mt-1 mb-[6px]'} to={`/artist/${artist.id}`}>
                    <img src={imageUrl(artist.image, "small")} alt={artist.name} />
                    <div>
                        <span style={{fontSize: 14}}>{artist.name}</span>
                        <div className="color-fg-secondary">{artist.listen_count} plays</div>
                    </div>
                    </Link>
                </div>
            );
        }
    }
}
