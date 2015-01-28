(function () {
    angular
        .module('gohub')
        .controller('reposController', reposController);

    function reposController($scope, reposService) {
        $scope.repositories = reposService.getRepositories();
        $scope.addRepository = function(repo) {
            reposService.addRepository(repo)
        }
    }
})();