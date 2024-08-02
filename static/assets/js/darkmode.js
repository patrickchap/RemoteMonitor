var themeToggleDarkIcon = document.getElementById('theme-toggle-dark-icon');
var themeToggleLightIcon = document.getElementById('theme-toggle-light-icon');

// Change the icons inside the button based on previous settings
if (localStorage.theme === 'dark' || (!('theme' in localStorage))) {
  document.documentElement.classList.add('dark')
  themeToggleLightIcon.classList.remove('hidden');
} else {
  document.documentElement.classList.remove('dark')
  themeToggleDarkIcon.classList.remove('hidden');
}


var themeToggleBtn = document.getElementById('theme-toggle');

themeToggleBtn.addEventListener('click', function() {

  // toggle icons inside button
  themeToggleDarkIcon.classList.toggle('hidden');
  themeToggleLightIcon.classList.toggle('hidden');

  // if set via local storage previously
  if (localStorage.theme == 'dark' || localStorage.theme == 'light') {
    if (localStorage.theme === 'light') {
      console.log("set dark");
      document.documentElement.classList.add('dark')
      localStorage.theme = 'dark';
    } else {

      console.log("set light");
      document.documentElement.classList.remove('dark');
      localStorage.theme = 'light';
    }
    // if NOT set via local storage previously
  } else {
    if (document.documentElement.classList.contains('dark')) {
      document.documentElement.classList.remove('dark');
      localStorage.theme = 'light';
    } else {
      document.documentElement.classList.add('dark');
      localStorage.theme = 'dark';
    }
  }

});
