package field

type StructField struct {
	keys   []string
	values []Field
}

func MallocStructField(cap int) *StructField {
	return &StructField{
		keys:   make([]string, 0, cap),
		values: make([]Field, 0, cap),
	}
}

func (f *StructField) Set(key string, value Field) {
	f.keys = append(f.keys, key)
	f.values = append(f.values, value)
}

func (f *StructField) Length() int {
	return len(f.keys)
}

// At (i int) return  a readonly Field in location i
func (f *StructField) At(i int) (string, Field, error) {
	if i >= f.Length() {
		return "", nil, OverIndexError
	}

	return f.keys[i], f.values[i], nil
}
