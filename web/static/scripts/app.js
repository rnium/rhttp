window.onload = () => {
    const ui = SwaggerUIBundle({
        url: "/schema",
        dom_id: "#swagger-ui",
    });

    const imageUrl = "/static/images/icon.svg";

    const insertOnce = () => {
        const title = document.querySelector("#swagger-ui h1.title");
        if (title && !document.getElementById("inserted-gopher-img")) {
            const img = document.createElement("img");
            img.src = imageUrl;
            img.id = "logo-img";
            img.alt = "Rhttp";
            img.style.width = "48px";
            img.style.height = "48px";
            img.style.marginRight = "8px";
            title.parentNode.insertBefore(img, title);
            return true;
        }
        return false;
    };

    if (!insertOnce()) {
        const targetNode = document.getElementById("swagger-ui");
        if (targetNode) {
            const observer = new MutationObserver((mutations, obs) => {
                if (insertOnce()) obs.disconnect();
            });
            observer.observe(targetNode, { childList: true, subtree: true });
        }
    }
};
