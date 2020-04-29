package my_models

// Page - страница доступная шаблонизатору
type Page struct {
	Posts PostItemSlice
	Token string
}
type Post struct {
	Posts PostItem
	Token string
}
