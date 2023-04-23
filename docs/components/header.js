const headerTemplate = document.createElement('template');
headerTemplate.innerHTML = `
<style> 
@import url("https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css");
@import url("https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.3/font/bootstrap-icons.css");
</style>
<header>
    <nav class="navbar navbar-expand-lg navbar-light bg-light">
        <div class="container-fluid">
            <a class="navbar-brand nav-link active" href="index.html">Packit</a>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarSupportedContent">
                <ul class="navbar-nav me-auto mb-2 mb-lg-0">
                    <li class="nav-item">
                    <a class="nav-link" href="registry.html">Registry</a>
                    </li>
                    <li class="nav-item">
                    <a class="nav-link" href="upload.html">Upload</a>
                    </li>
                    <li class="nav-item">
                    <a class="nav-link" href="manage.html">Management</a>
                    </li>
                    <li class="nav-item">
                    <a class="nav-link" href="404lol">404</a>
                    </li>
                </ul>
                <div>
                    <button type="button" class="btn btn-primary" data-toggle="modal" data-target="#loginModal"> Login</button>
                    <button class="btn btn-outline-primary onclick="location.href='settings.html'" type="submit">Settings</button> 
                </div>
            </div>
        </div>
    </nav>
</header>
`
class Header extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        const shadowRoot = this.attachShadow({ mode: 'closed' });
    
        shadowRoot.appendChild(headerTemplate.content);
    }
}

customElements.define('header-component', Header);