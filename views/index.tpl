<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cats</title>
    <link rel="stylesheet" href="/static/css/styles.css">
</head>
<body>
    <div class="card">
        <div class="menu">
            <span class="menu-item" id="voting-menu">Voting</span>
            <span class="menu-item" id="breeds-menu">Breeds</span>
            <span class="menu-item" id="favs-menu">Favs</span>
        </div>
        <div class="content">
            <div id="voting-content">
                <img id="catImage" src="" alt="Random cat" class="cat-image">
                <div class="actions">
                    <button class="action-button" id="likeBtn" onclick="handleFavorite()">❤️</button>
                    <div>
                        <button class="action-button" id="upvoteBtn" onclick="submitVote(1)">👍</button>
                        <button class="action-button" id="downvoteBtn" onclick="submitVote(-1)">👎</button>
                    </div>
                </div>
                <div class="cat-info">
                    <p>Breed: <span id="catBreed">Loading...</span></p> 
                </div>
            </div>
            <div id="breeds-content">
                <div class="select-container">
                    <select id="breedSelect" class="search-bar"></select>
                </div>
                <div class="slider-container">
                    <div id="sliderImagesContainer">
                        <!-- Images will be dynamically added here -->
                    </div>
                    <div class="nav-dots" id="sliderNavDots"></div>
                </div>                     
                <div class="breed-info">
                    <span class="breed-title" id="breedName">Name</span>
                    <span class="breed-origin id-italic" id="breedId">id</span>
                    <p class="breed-title breed-origin"><span id="breedOrigin">Origin:</span></p>
                    <p class="breed-origin" id="breedDescription"></p>
                    <a class="wikipedia" id="wikiLink" target="_blank">WIKIPEDIA</a>
                </div>
            </div>
            <div id="favs-content">
                <div class="view-toggle">
                    <button class="view-toggle-item" id="grid-view">☷</button>
                    <button class="view-toggle-item" id="list-view">☰</button>
                </div>
                <div id="favoritesSection" style="display: none;">
                    <div class="favorites-grid" id="favoritesGrid"></div>
                </div>
            </div>
        </div>
    </div>

    <script src="/static/js/scripts.js"></script>
</body>
</html>