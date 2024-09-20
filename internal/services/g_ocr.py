from flask import Flask, request, jsonify
import io
import requests
import re
from google.cloud import vision

app = Flask(__name__)

client = vision.ImageAnnotatorClient()

def fetch_image_from_url(image_url):
    response = requests.get(image_url)
    if response.status_code != 200:
        raise Exception(f"Failed to fetch image from {image_url}, status code: {response.status_code}")
    
    return io.BytesIO(response.content)

def detect_text_from_image(image_bytes):
    image = vision.Image(content=image_bytes.getvalue())
    response = client.text_detection(image=image)
    texts = response.text_annotations

    if not texts:
        return "No text detected."

    full_text = texts[0].description
    pattern = r'(o{3,}|O{3,}|0{3,})'
    cleaned_text = re.sub(pattern, '\n', full_text)

    return cleaned_text

@app.route("/detect-text", methods=["POST"])
def detect_text():
    data = request.get_json()
    if "image_url" not in data:
        return jsonify({"error": "No image URL provided."}), 400
    
    image_url = data["image_url"]
    
    try:
        image_bytes = fetch_image_from_url(image_url)
        
        detected_text = detect_text_from_image(image_bytes)
        
        return jsonify({"text": detected_text})
    
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
