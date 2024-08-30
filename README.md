# Beego Cat API 

A web application that allows users to explore different cat breeds, view images, and vote on their favorite cats. This project is built using the Beego web framework for the backend and vanilla JavaScript for the frontend. In this project, a simple API has been built using the Beego framework in Go. The API provides functionalities to interact with the [The Cat API](https://thecatapi.com/), including fetching random cat images, voting on cats, managing favorites and viewing cats by breeds.

## Features

- Fetches random cats
- Browse a list of cat breeds
- View detailed information about each breed
- See images of cats from selected breeds
- Vote (like/dislike) or add favorite feature on cat images
- Simple and intuitive user interface

## Prerequisites

- Go 1.16+
- Beego Framework v2
- An API key from The Cat API
  
## Project Structure

  ```plaintext
   go/
    └── src/
         └── cat_api/
              ├── conf/
              │   └── app.conf           # Configuration file for the application
              ├── controllers/
              │   └── cat_controller.go  # Controller logic (e.g., handling HTTP requests)
              ├── main.go                # Main application entry point
              ├── routers/
              │   └── router.go          # Routes definition and setup
              ├── static/
              │   ├── css/               # CSS files
              │   └── js/                # JavaScript files
              ├── tests/
              │   └── default_test.go    # Test files
              └── views/
                  └── index.tpl          # HTML template files for views
  ```
  
## Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/nafiahossain/cat_api.git
   cd cat_api
   ```

2. **Install Beego v2 and Bee v2 tool**
 
   Check go version

   ```bash
    go version
   ```

   If go exists in your system then install the Beego Framework using the following commands:

   ```bash
    go install github.com/beego/beego/v2@latest
    go install github.com/beego/bee/v2@latest
   ```

   You can check the Bee version for installtion confirmation:

   ```bash
   bee version
   ```

3. **Install dependencies**

   ```bash
   go mod tidy
   ```

4. **Set Up The Cat API Key**

    1. Get your API key from [The Cat API](https://thecatapi.com/#pricing).
    2. Add your API key to the conf/app.conf file:

        ```ini
        cat_api_key = your_actual_api_key_from_thecatapi.com
        ```

5. **Configurations**

    The application settings are managed in the conf/app.conf file.  Rename the app.conf.sample file to app.conf and then add you the Cat API key. You can configure the following:

    - httpport: The port on which the application will run (default is 8080).
    - cat_api_key: Your API key for The Cat API.

6. **Run the Project**

    To start the Beego application, run:

    ```bash
    bee run
    ```

    The application will start at http://localhost:8080.

## API Endpoints

Here are some of the API endpoints available in this project:

- GET /api/cat: Fetches a random cat image. Each contains the url for the image file, along with its width, height and breed information (if available).
- POST /api/vote: Submit a vote for a cat image.
- POST /api/favorite: Adds a cat image to favorites.
- GET /api/favorites: Retrieves the list of favorite cat images.
- GET /api/breeds: Fetch list of cat breeds.
- GET /api/breed/:id: Fetch information about a specific breed.

## Usage

- Upon loading, the application will display a cat image with its breed on `Voting`.
- Use the thumbs up/down buttons to vote on the current cat image.
- Use the Heart button to add the current cat image to the favorites list that can be seen in the `Favs` Section.
- In the `Breeds` section, the application will display information about the first cat breed in the list.
- Use the dropdown menu to select different cat breeds.

## Static Files

Static files such as CSS and JavaScript are located in the static directory. These files are automatically served by Beego when referenced in your HTML templates.

- CSS: Located in static/css/.
- JavaScript: Located in static/js/.

## Templates

- HTML templates are stored in the views/ directory. The main template used in this project is:

    - index.tpl: The primary template file rendering the UI for cat images and interactions.

## Contributing

Contributions are welcome! Please fork the repository and create a pull request to add features or fix issues.
