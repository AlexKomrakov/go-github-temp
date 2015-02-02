(function () {
    angular
        .module('gohub')
        .service('reposService', reposService);

    function reposService($resource) {
        var reposService = this;
        var repositories = { data: [] };
        reposService.api = $resource("/repos");
        reposService.loadRepositories = function () {
            reposService.api.query(function(result){
                repositories.data = result;
            })
        };
        reposService.getRepositories = function () {
            reposService.loadRepositories();
            return repositories;
        };
        reposService.addRepository = function (repo) {
            reposService.api.save(repo, function(){
                reposService.loadRepositories();
            })
        };
    }
})();