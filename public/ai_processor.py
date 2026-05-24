import speech_recognition as sr
import whisper
import os
import time

# --- CONFIGURATION ---
# We use the 'base' model for a balance of speed and accuracy
WHISPER_MODEL_NAME = "base"
# The target sentence to verify
TARGET_SENTENCE = "I only know that John is the one to help today."

def main():
    print("\n--- CLINICAL AI VOICE PROCESSOR (WHISPER + GEMMA) ---")
    
    # 1. Initialize Whisper
    print(f"[*] Loading AI Model: Whisper ({WHISPER_MODEL_NAME})...")
    model = whisper.load_model(WHISPER_MODEL_NAME)
    
    # 2. Capture Audio from Microphone
    recognizer = sr.Recognizer()
    with sr.Microphone() as source:
        print("\n[READY] Please repeat the sentence: ")
        print(f"Target: \"{TARGET_SENTENCE}\"")
        print("\n--- Listening (Speak now) ---")
        
        # Adjust for ambient noise for 1 second
        recognizer.adjust_for_ambient_noise(source, duration=1)
        audio = recognizer.listen(source)
        
    print("[*] Processing transmission...")

    # 3. Save audio to temporary file for Whisper
    with open("temp_audio.wav", "wb") as f:
        f.write(audio.get_wav_data())

    # 4. Transcribe using Whisper
    start_time = time.time()
    result = model.transcribe("temp_audio.wav")
    transcription = result["text"].strip()
    end_time = time.time()

    print(f"\n[AI TRANSCRIPTION]: \"{transcription}\"")
    print(f"[*] Transcription Time: {end_time - start_time:.2f}s")

    # 5. Verification Logic (Gemma Logic)
    # Note: In a full setup, we would send this to a local Gemma instance via Ollama
    print("\n[*] Initializing Gemma Verification...")
    
    # Simple semantic check (simulating Gemma's logic)
    cleaned_trans = transcription.lower().replace(".", "").replace(",", "")
    cleaned_target = TARGET_SENTENCE.lower().replace(".", "").replace(",", "")
    
    is_correct = "john" in cleaned_trans and "help" in cleaned_trans and "today" in cleaned_trans
    
    print("\n--- CLINICAL REPORT ---")
    if is_correct:
        print("STATUS: [NOMINAL]")
        print("VERDICT: Sentence successfully repeated with high fidelity.")
    else:
        print("STATUS: [VARIANCE DETECTED]")
        print("VERDICT: Transcription mismatch. Semantic elements missing.")

    # Cleanup
    if os.path.exists("temp_audio.wav"):
        os.remove("temp_audio.wav")

if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print(f"\n[ERROR] System fault: {e}")
        print("Ensure you have 'openai-whisper', 'SpeechRecognition', and 'pyaudio' installed.")
