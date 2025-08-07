# MdConv

This project aims to convert markdown to pdf, going via transitioning to HTML.

## How to run this code?

- Clone the repo -
    ```bash
    git clone https://github.com/architmishra-15/MdConv.git md-to-pdf
    cd md-to-pdf
    ```

- Install the dependencies -
    ```bash
    go mod tidy
    ```

- Run the program and pass the name of markdown
    ```bash
    go run . test.md
    ```

> Todo: Correct the syantax highlighting to make it confide till the code block, and add the final html -> pdf.
