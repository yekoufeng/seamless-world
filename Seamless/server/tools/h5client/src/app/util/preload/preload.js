function preloadProgress(message) {
    var element = document.getElementById('page-preload-progress');
    element.innerHTML = message;
}

function closePreload(app) {
    app.run(['$animate', '$timeout', function Execute($animate, $timeout) {
        $timeout(function () {
            var element = angular.element(document.getElementById('page-preload'));
            $animate.addClass(element, 'preload-fade').then(function () {
                element.remove();
            });
        }, 500);
    }]);
}
