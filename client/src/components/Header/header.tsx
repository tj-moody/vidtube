function Header() {
    return (
        <header
            style={{
                background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
                padding: "1rem 2rem",
                boxShadow: "0 2px 10px rgba(0, 0, 0, 0.1)",
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
            }}
        >
            <h1
                style={{
                    color: "white",
                    margin: 0,
                    fontSize: "1.5rem",
                    fontWeight: "bold",
                }}
            >
                VidTube
            </h1>
        </header>
    );
}

export default Header;
