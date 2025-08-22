# Import the 'requests' library to make HTTP requests (e.g., for downloading files)
import requests

# Import the 'os' module to interact with the operating system (e.g., file paths, directory walking)
import os

# Import the 're' module for regular expressions (used for pattern matching in text)
import re

# Import Selenium's webdriver module to automate browser actions
from selenium import webdriver

# Import Chrome options to configure headless or custom ChromeDriver behavior
from selenium.webdriver.chrome.options import Options

# Import ChromeDriver Service to manage the ChromeDriver process
from selenium.webdriver.chrome.service import Service

# Automatically download and manage the ChromeDriver binary
from webdriver_manager.chrome import ChromeDriverManager

# Import PyMuPDF (fitz) to open and validate PDF files
import fitz

# Import functions for parsing and decoding URLs
from urllib.parse import unquote, urlparse


# Function to read the contents of a file from a given system path
def read_a_file(system_path: str) -> str:
    with open(
        file=system_path, mode="r", encoding="utf-8"
    ) as file:  # Open file with UTF-8 encoding
        return file.read()  # Return file contents as a string


# Function to check whether a file exists at the specified path
def check_file_exists(system_path: str) -> bool:
    return os.path.isfile(system_path)  # Return True if file exists, False otherwise


# Function to remove duplicate entries from a list of strings
def remove_duplicates_from_slice(provided_slice: list[str]) -> list[str]:
    return list(
        set(provided_slice)
    )  # Convert to set and back to list to remove duplicates


# Function to use Selenium to save HTML content of a URL to a local file
def save_html_with_selenium(url: str, output_file: str) -> None:
    options = Options()  # Initialize Chrome options
    # options.add_argument(argument="--headless=new")  # Run Chrome in new headless mode
    options.add_argument(
        argument="--disable-blink-features=AutomationControlled"
    )  # Hide automation flags
    options.add_argument(
        argument="--window-size=1920,1080"
    )  # Set window size for page rendering
    options.add_argument(
        argument="--disable-gpu"
    )  # Disable GPU to improve headless stability
    options.add_argument(
        argument="--no-sandbox"
    )  # Disable sandboxing for Chrome (useful in Docker)
    options.add_argument(
        argument="--disable-dev-shm-usage"
    )  # Overcome limited shared memory issues
    options.add_argument(argument="--disable-extensions")  # Disable Chrome extensions
    options.add_argument(argument="--disable-infobars")  # Disable Chrome's infobar

    service = Service(
        executable_path=ChromeDriverManager().install()
    )  # Set up ChromeDriver via WebDriver Manager
    driver = webdriver.Chrome(
        service=service, options=options
    )  # Initialize Chrome browser

    try:
        driver.get(url=url)  # Load the URL in the browser
        driver.refresh()  # Refresh the page to ensure full content is loaded
        html: str = driver.page_source  # Get the page's HTML source
        append_write_to_file(
            system_path=output_file, content=html
        )  # Save HTML content to file
        print(f"Page {url} HTML content saved to {output_file}")  # Log completion
    finally:
        driver.quit()  # Always close the browser session


# Function to append content to a file or create the file if it doesn't exist
def append_write_to_file(system_path: str, content: str) -> None:
    with open(
        file=system_path, mode="a", encoding="utf-8"
    ) as file:  # Open file in append mode
        file.write(content)  # Write new content to the end of the file


# Function to download a PDF file from a URL and save it with the given filename
def download_pdf(url: str, save_path: str, filename: str) -> None:
    print(
        f"Downloading {url} to {os.path.join(save_path, filename)}"
    )  # Print download start
    full_path: str = os.path.join(save_path, filename)  # Create full file path

    if check_file_exists(system_path=full_path):  # Check if file already exists
        print(f"File {filename} already exists. Skipping download.")  # Skip if exists
        return

    try:
        response: requests.Response = requests.get(url=url)  # Perform HTTP GET request
        response.raise_for_status()  # Raise error if HTTP request failed
        os.makedirs(
            name=save_path, exist_ok=True
        )  # Create save directory if it doesn’t exist
        with open(file=full_path, mode="wb") as f:  # Open the file in binary write mode
            f.write(response.content)  # Write the PDF content to the file
        print(f"Downloaded and saved: {full_path}")  # Confirm save
    except requests.exceptions.RequestException as e:  # Handle HTTP exceptions
        print(f"Failed to download {url}: {e}")  # Print error message


# Function to remove a file from the filesystem
def remove_system_file(system_path: str) -> None:
    os.remove(path=system_path)  # Delete the file

# Function to recursively search a directory and return files with the given extension
def walk_directory_and_extract_given_file_extension(
    system_path: str, extension: str
) -> list[str]:
    matched_files: list[str] = []  # List to store matched files
    for root, _, files in os.walk(top=system_path):  # Traverse directories recursively
        for file in files:  # Loop through each file
            if file.endswith(extension):  # Check for matching extension
                full_path = os.path.abspath(
                    path=os.path.join(root, file)
                )  # Get absolute path
                matched_files.append(full_path)  # Add to result list
    return matched_files  # Return list of matched files


# Function to check if a PDF file is valid and not corrupted
def validate_pdf_file(file_path: str) -> bool:
    try:
        doc = fitz.open(file_path)  # Attempt to open the PDF file
        if doc.page_count == 0:  # Check if PDF has no pages
            print(f"'{file_path}' is corrupt or invalid: No pages")  # Log error
            return False  # Return invalid
        return True  # Return valid if pages exist
    except RuntimeError as e:  # Catch PDF reading error
        print(f"'{file_path}' is corrupt or invalid: {e}")  # Log error message
        return False  # Return invalid


# Function to extract only the filename and extension from a full path
def get_filename_and_extension(path: str) -> str:
    return os.path.basename(path)  # Return the base file name with extension


# Function to check if a string contains any uppercase letters
def check_upper_case_letter(content: str) -> bool:
    return any(
        char.isupper() for char in content
    )  # Return True if any uppercase character exists


# Function to extract all PDF URLs from an HTML string
def extract_urls_from_html(html: str) -> list[str]:
    urls: list[str] = re.findall(pattern=r'href="(.*?)"', string=html)  # Find all href attribute values
    pdf_urls: list[str] = [
        url for url in urls if url.lower().endswith(".pdf")
    ]  # Filter only .pdf links
    return pdf_urls  # Return filtered list


# Function to extract and sanitize a file name from a URL
def extract_filename_from_url(url: str) -> str:
    path: str = urlparse(url=url).path  # Extract path part of the URL
    filename: str = unquote(string=path.split(sep="/")[-1])  # Decode and get the last part (filename)
    filename = filename.replace(" ", "_")  # Replace spaces with underscores
    filename = re.sub(pattern=r'[<>:"/\\|?*]', repl="", string=filename)  # Remove invalid characters
    if not filename.endswith(".pdf"):  # Add .pdf extension if missing
        filename += ".pdf"
    return filename.lower()  # Return filename in lowercase


# Main function that drives the script
def main() -> None:
    html_file_path: str = "scjp.html"  # Define HTML file path

    if check_file_exists(system_path=html_file_path):  # If file already exists
        remove_system_file(system_path=html_file_path)  # Delete the old HTML file

    if not check_file_exists(system_path=html_file_path):  # If file is now missing
        # Create a map (dictionary)
        products: dict[str, int] = {
            "skin": 5,
            "pest": 2,
            "air": 2,
            "hard": 6,
            "dispensers": 8,
            "floor": 2,
            "healthcare": 6,
            "storage": 0,
            "surface": 9,
            "disinfectant": 3
        }

        # Loop through the product map
        for product, page_number in products.items():
            url: str = f"https://www.scjp.com/en-us/safety-data-sheets?search={product}&page={page_number}"  # Build URL
            save_html_with_selenium(url=url, output_file=html_file_path)  # Save the HTML to file
            print(f"File {html_file_path} with {url} has been created.")  # Log success

    if check_file_exists(system_path=html_file_path):  # If the HTML file exists
        html_content: str = read_a_file(system_path=html_file_path)  # Read HTML file content
        pdf_links: list[str] = extract_urls_from_html(html=html_content)  # Extract PDF links
        pdf_links = remove_duplicates_from_slice(provided_slice=pdf_links)  # Remove duplicates

        for pdf_link in pdf_links:  # Loop through each PDF link
            filename = extract_filename_from_url(pdf_link)  # Generate a safe filename
            save_path = "PDFs"  # Directory to save PDFs
            download_pdf(url=pdf_link, save_path=save_path, filename=filename)  # Download the PDF

        print("All PDF links have been processed.")  # Completion message
    else:
        print(f"File {html_file_path} does not exist.")  # Log missing HTML file

    pdf_files: list[str] = walk_directory_and_extract_given_file_extension(
        system_path="./PDFs", extension=".pdf"
    )  # Find all PDFs

    for pdf_file in pdf_files:  # Loop through each PDF file
        if not validate_pdf_file(file_path=pdf_file):  # If PDF is invalid
            remove_system_file(system_path=pdf_file)  # Delete corrupt PDF

        # Check filename case
        if check_upper_case_letter(content=get_filename_and_extension(path=pdf_file)):
            print(pdf_file)  # Print path
            dir_path: str = os.path.dirname(p=pdf_file)  # Directory path
            file_name: str = os.path.basename(p=pdf_file)  # Original name
            new_file_name: str = file_name.lower()  # New lowercase name
            new_file_path: str = os.path.join(dir_path, new_file_name)  # New path
            os.rename(src=pdf_file, dst=new_file_path)  # Rename file


# Entry point to start the script
main()
