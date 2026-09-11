import { useRef, useState } from "react";
import styles from "./videoupload.module.css";

type VideoUploadProps = {
    onVideoAdded: () => void;
};

const baseUrl = "http://localhost:8080/";

const VideoUpload = ({ onVideoAdded }: VideoUploadProps) => {
    const [title, setTitle] = useState("");
    const [authorID, setAuthorID] = useState("");
    const [file, setFile] = useState<File | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    async function handleSubmit(e: React.SubmitEvent) {
        e.preventDefault();

        if (isSubmitting) return;

        if (!title.trim()) {
            alert("Please enter a video title");
            return;
        }
        const parsedAuthorID = Number(authorID);
        if (
            !authorID.trim() ||
            !Number.isInteger(parsedAuthorID) ||
            parsedAuthorID <= 0
        ) {
            alert("Please enter a valid author ID");
            return;
        }

        if (
            !file ||
            file.size === 0 ||
            !file.name.toLowerCase().endsWith(".mp4")
        ) {
            alert("Please select a non-empty MP4 video");
            return;
        }

        setIsSubmitting(true);

        try {
            const response = await fetch(baseUrl + "api/v1/videos", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    title: title.trim(),
                    authorID: parsedAuthorID,
                }),
            });

            if (!response.ok) {
                const errorData = await response.text();
                throw new Error(
                    errorData || `HTTP error! status: ${response.status}`,
                );
            }

            const data = await response.json();
            if (typeof data.uploadURL !== "string" || !data.uploadURL) {
                throw new Error("Server did not return an upload URL");
            }

            const uploadResponse = await fetch(data.uploadURL, {
                method: "PUT",
                headers: { "Content-Type": "video/mp4" },
                body: file,
            });

            if (!uploadResponse.ok) {
                throw new Error(
                    `Video upload failed (HTTP ${uploadResponse.status})`,
                );
            }

            alert(
                "Video uploaded! It will appear once the server finishes processing the upload.",
            );
            setTitle("");
            setAuthorID("");
            setFile(null);
            if (fileInputRef.current) fileInputRef.current.value = "";
            onVideoAdded();
        } catch (error) {
            console.error("Error uploading video:", error);
            alert(
                `Failed to upload video: ${error instanceof Error ? error.message : "Unknown error"}`,
            );
        } finally {
            setIsSubmitting(false);
        }
    }

    return (
        <div className={styles.upload_container}>
            <h2>Upload Video</h2>
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
                <div className={styles.form_group}>
                    <label htmlFor="videoFile">MP4 Video *</label>
                    <input
                        ref={fileInputRef}
                        id="videoFile"
                        type="file"
                        accept=".mp4,video/mp4"
                        onChange={(e) => setFile(e.target.files?.[0] ?? null)}
                        className={styles.input}
                        disabled={isSubmitting}
                    />
                </div>
                <button
                    type="submit"
                    className={styles.submit_button}
                    disabled={isSubmitting}
                >
                    {isSubmitting ? "Uploading..." : "Upload Video"}
                </button>
            </form>
        </div>
    );
};

export default VideoUpload;
