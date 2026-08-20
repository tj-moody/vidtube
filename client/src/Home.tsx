import { useState, useEffect } from "react";
import Header from "./components/Header/header";
import VideoPane from "./components/VideoPane/videopane";
import VideoUpload from "./components/VideoUpload/videoupload";
import styles from "./vidtube.module.css";

export type Video = {
    id: string;
    title: string;
    author: string;
    views: string;
};

const Vidtube = () => {
    const [videos, setVideos] = useState<Video[]>([]);

    const loadVideos = async () => {
        const url = "http://localhost:8080/";

        const data = await fetch(url + "api/v1/videos" + "?page=1&count=10")
            .then((response) => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .catch((error) => {
                console.error(`Error fetching videos at ${url}: `, error);
                alert(`Error fetching videos at ${url}`);
                return null;
            });
        if (data == null) {
            return;
        }
        setVideos(data);
    };

    useEffect(() => {
        loadVideos();
    }, []);

    return (
        <>
            <Header />
            <VideoUpload onVideoAdded={() => loadVideos()} />
            <div className={styles.video_grid}>
                {videos.map((video) => (
                    <VideoPane key={video.id} video={video} />
                ))}
            </div>
        </>
    );
};

export default Vidtube;
