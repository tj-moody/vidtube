import type { Video } from "../../Home.tsx";
import styles from "./videopane.module.css"

type Props = {
    video: Video
}
const VideoPane = ({ video }: Props) => {
    const PLACEHOLDER_THUMBNAIL_URL = "https://placehold.co/160x90";
    console.log(video)
    return (
        <div className={styles.video_pane}>
            <a href={`/watch?v=${encodeURIComponent(video.id)}`} target="_blank">
                <img src={PLACEHOLDER_THUMBNAIL_URL} alt={video.title} />
            </a>
            <div key={video.id}>
                <h3>{video.title}</h3>
                <p>
                    {video.author} — {video.views} views
                </p>
            </div>
        </div>
    );
};

export default VideoPane;
