package main

import (
	"ascii-art-web/ascii" // Import your local package
	"html/template"
	"net/http"
	"os"
)

// PageData struct to hold data sent to HTML
// هيكل البيانات الذي سنرسله لملف الـ HTML
type PageData struct {
	Result string
}

func main() {
	// 1. Handle Routes
	// تحديد الروابط (Endpoints)
	http.HandleFunc("/", homeHandler)           // الصفحة الرئيسية
	http.HandleFunc("/ascii-art", asciiHandler) // صفحة المعالجة

	// 2. Start Server
	// تشغيل السيرفر على المنفذ 8080
	println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		println("Error starting server:", err)
	}
}

// GET / : Handler for the main page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// 404 Check: If the path is not exactly "/", return 404
	// التأكد أن الرابط هو الصفحة الرئيسية فقط
	if r.URL.Path != "/" {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}

	// Render the template
	renderTemplate(w, nil) // Send nil because no data yet
}

// POST /ascii-art : Handler to process the form
func asciiHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Check Method: Must be POST
	// التأكد أن الطلب من نوع POST
	if r.Method != http.MethodPost {
		http.Error(w, "400 - Method Not Allowed", http.StatusBadRequest)
		return
	}

	// 2. Get Form Values
	// استخراج البيانات من الفورم
	text := r.FormValue("text")
	banner := r.FormValue("banner")

	// Validation: Check if banner is valid
	if banner != "standard" && banner != "shadow" && banner != "thinkertoy" {
		http.Error(w, "400 - Bad Request (Invalid Banner)", http.StatusBadRequest)
		return
	}

	// Check if banners exist
	if _, err := os.Stat("banners/" + banner + ".txt"); os.IsNotExist(err) {
		http.Error(w, "404 - Banner File Not Found", http.StatusNotFound)
		return
	}

	// 3. Generate ASCII Art
	// استدعاء دالة الرسم التي أنشأناها
	art, err := ascii.GenerateArt(text, banner)
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 4. Send Data back to HTML
	// إرسال النتيجة لصفحة الويب
	data := PageData{
		Result: art,
	}
	renderTemplate(w, data)
}

// Helper function to parse and execute template
// دالة مساعدة لقراءة ملف HTML
func renderTemplate(w http.ResponseWriter, data interface{}) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "500 - Template Error", http.StatusInternalServerError)
		return
	}
	// 200 OK is implicit here
	tmpl.Execute(w, data)
}
