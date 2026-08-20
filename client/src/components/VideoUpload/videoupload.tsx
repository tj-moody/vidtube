import { useState } from "react";
import styles from "./videoupload.module.css";

type VideoUploadProps = {
    onVideoAdded: () => void;
};

const baseUrl = "http://localhost:8080/";

const VideoUpload = ({ onVideoAdded }: VideoUploadProps) => {
    const [title, setTitle] = useState("");
    const [authorID, setAuthorID] = useState("");
    const [isSubmitting, setIsSubmitting] = useState(false);

    async function handleSubmit(e: React.SubmitEvent) {
    e.preventDefault();

    if (!title.trim()) {
        alert("Please enter a video title");
        return;
    }
    const parsedAuthorID = Number(authorID);
    if (!authorID.trim() || !Number.isInteger(parsedAuthorID) || parsedAuthorID <= 0) {
        alert("Please enter a valid author ID");
        return;
    }

    setIsSubmitting(true);

    try {
        const response = await fetch(baseUrl + "api/v1/videos", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ title, authorID: parsedAuthorID }),
        });

        if (!response.ok) {
            const errorData = await response.text();
            throw new Error(
                errorData || `HTTP error! status: ${response.status}`,
            );
        }

        const data = await response.json();
        console.log("Got presigned upload URL:", data.uploadURL);

        alert("Video added successfully!");
        setTitle("");
        setAuthorID("");
        onVideoAdded();
    } catch (error) {
        console.error("Error adding video:", error);
        alert(`Failed to add video: ${error}`);
    } finally {
        setIsSubmitting(false);
    }
}

    return (
        <div className={styles.upload_container}>
            <h2>Add Video</h2>
            <form onSubmit={handleSubmit} className={styles.upload_form}>
                <div className={styles.form_group}>
                    <label htmlFor="title">Video Title *</label>
                    <input
                        id="title"
                        type="text"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        placeholder="Enter video title"
                        className={styles.input}
                        disabled={isSubmitting}
                    />
                </div>
                <div className={styles.form_group}>
                    <label htmlFor="authorID">Author ID *</label>
                    <input
                        id="authorID"
                        type="number"
                        min="1"
                        value={authorID}
                        onChange={(e) => setAuthorID(e.target.value)}
                        placeholder="Enter author user ID"
                        className={styles.input}
                        disabled={isSubmitting}
                    />
                </div>
                <button
                    type="submit"
                    className={styles.submit_button}
                    disabled={isSubmitting}
                >
                    {isSubmitting ? "Adding..." : "Add Video"}
                </button>
            </form>
        </div>
    );
};

export default VideoUpload;
