const loginTemplate = document.createElement('template');
// RELEVANT JS TO USE THIS MODAL
// $('#myModal').on('shown.bs.modal', function () {
//     $('#myInput').trigger('focus')
//   })
loginTemplate.innerHTML = `
<!-- Modal -->
<div class="modal fade" id="loginModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <div class="col-md-6">
                    <h5 class="modal-title" id="loginModalTitle1">Sign Up</h5>
                </div>
                <div class="col-md-6">
                    <h5 class="modal-title" id="loginModalTitle2">Login</h5>
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true">&times;</span>
                    </button>
                </div>
            </div>
            <div class="modal-body">
                <div class="col-md-6">
                    something 1
                </div>
                <div class="col-md-6">
                    something 2
                </div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
            </div>
        </div>
    </div>
</div>
`;

class LoginModal extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        const shadowRoot = this.attachShadow({ mode: 'closed' });

        shadowRoot.appendChild(loginTemplate.content);
    }
}

customElements.define('login-modal-component', LoginModal);