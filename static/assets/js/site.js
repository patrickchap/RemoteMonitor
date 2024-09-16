
console.log('site.js loaded');

document.body.addEventListener('htmx:afterOnLoad', (evt) => {
  if (evt.detail.target.id === 'hostservices') {
    console.log('modal is loaded');

    const $targetEl = document.getElementById('crud-modal');
    console.log(evt.detail.target);

    if ($targetEl) {
      console.log('modal is loaded');
      const options = {
        placement: 'center',
        backdrop: 'dynamic',
        backdropClasses: 'bg-gray-900/50 dark:bg-gray-900/80 fixed inset-0 z-40',
        closable: true,
        onHide: () => {
          console.log('modal is hidden');
        },
        onShow: () => {
          console.log('modal is shown');
        },
        onToggle: () => {
          console.log('modal has been toggled');
        }
      };


      /*
      * $targetEl: required
      * options: optional
      */
      const modal = new Modal($targetEl, options);

      const addServiceButton = document.getElementById("add-service");
      addServiceButton.addEventListener('click', () => {
        modal.show();
      });

      const closeButton = document.getElementById("close-model");
      closeButton.addEventListener('click', () => {
        modal.hide();
      });
    }
  }
  // o}ptions with default values

});




