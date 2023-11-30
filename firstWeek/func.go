package main

/*
有一个问题：
按我这种方式的话，没有缩容的机会
我查看地址，每次使用的slice地址是一致的，但我不知道怎么获取起始slice的容量，
我的想法是获取起始slice的容量，然后获取现在的slice的长度，
1.如果容量<256的话,长度 < 容量 * 0.5  开始缩容
2.如果容量>256的话,长度 < 容量 * 0.75 开始缩容
*/
func SliceDelete[T any](index int, vals []T) []T {
	vals_length := len(vals)

	if index < 0 || index >= vals_length {
		panic("index 不合法")
	}

	vals_cap := cap(vals)
	length_half := vals_length / 2
	if vals_length < (vals_cap / 2) {
		new_vals := make([]T, vals_length)
		for tmpIndex, value := range vals {
			if tmpIndex != index {
				new_vals = append(new_vals, value)
			}
		}
		return new_vals
	} else {
		i := index
		if index <= length_half {
			for i = index; i > 0; i-- {
				vals[i] = vals[i-1]
			}
			return vals[i+1:]
		} else {
			for i = index; i < vals_length-1; i++ {
				vals[i] = vals[i+1]
			}
			return vals[:i]
		}
	}
}
