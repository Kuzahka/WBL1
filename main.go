package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

// 1 Дана структура Human (с произвольным набором полей и методов).
// Реализовать встраивание методов в структуре Action от родительской
//  структуры Human (аналог наследования).
	type Human struct {
		Name string
		Age  int
	}
	func (h Human) GetName() string {
		return h.Name
	}
	func (h *Human) SetName(name string) {
		h.Name = name
	}
	func (h Human) GetAge() int {
		return h.Age
	}
	type Action struct {
		Human
		ActionType string
	}
	func (a Action) GetActionType() string {
		return a.ActionType
	}
	//Action встраивает Human, поэтому все методы и поля Human доступны напрямую через Action
func main(){
	action := Action{
		Human: Human{
			Name: "John",
			Age:  30,
		},
		ActionType: "Running",
	}
	//Поля структуры Human доступны через Action, как если бы они были частью Action
	// Используем методы Human
	fmt.Println("Name:", action.GetName())
	fmt.Println("Age:", action.GetAge())

	// Используем метод Action
	fmt.Println("Action Type:", action.GetActionType())

//		// Изменяем имя через метод Human
//		action.SetName("Mike")
//		fmt.Println("Updated Name:", action.GetName())
//	}
//
// ////////////////////////////////////////////
// ////////////////////////////////////////////
// ///////////////////////////////////////////
// 2 Написать программу, которая конкурентно
//	рассчитает значение квадратов чисел взятых
//	из массива (2,4,6,8,10) и выведет их квадраты в stdout.
//
// 1 Решение
func main() {
	// Исходный массив чисел
	numbers := []int{2, 4, 6, 8, 10}

	// Группа ожидания для синхронизации горутин
	var wg sync.WaitGroup

	// Канал для передачи результатов
	results := make(chan int, len(numbers))

	// Запуск горутин для расчета квадратов
	for _, num := range numbers {
		wg.Add(1) // Увеличиваем счетчик группы ожидания
		go func(n int) {
			defer wg.Done()  // Уменьшаем счетчик при завершении
			results <- n * n // Отправляем результат в канал
		}(num)
	}

	// Закрытие канала после завершения всех горутин
	wg.Wait()      // Ждем завершения всех горутин
	close(results) // Закрываем канал

	// Чтение и вывод результатов
	for square := range results {
		fmt.Println(square)
	}

}

// ////////////////////////////////////////
// /////////////////////////////////////////
// 2 Решение
func main() {
	numbers := []int{2, 4, 6, 8, 10}
	squares := make([]int, len(numbers))

	var wg sync.WaitGroup

	for i, num := range numbers {
		wg.Add(1)
		go func(index, value int) {
			defer wg.Done()
			squares[index] = value * value
		}(i, num)
	}

	wg.Wait()

	for _, square := range squares {
		fmt.Println(square)
	}
}
// ////////////////////////////////////////
// /////////////////////////////////////////
// 3 Решение
func main() {
	numbers := []int{2, 4, 6, 8, 10}
	var wg sync.WaitGroup

	for _, num := range numbers {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fmt.Println(n * n)
		}(num)
	}

	wg.Wait()
}
// Дана последовательность чисел: 2,4,6,8,10. Найти сумму их квадратов(2^2+3^2+4^2….)
//  с использованием конкурентных вычислений.
func main() {
	// Исходная последовательность чисел
	numbers := []int{2, 4, 6, 8, 10}

	// Переменная для хранения суммы
	var sum int
	var mu sync.Mutex // Мьютекс для синхронизации доступа к переменной суммы

	// Группа ожидания
	var wg sync.WaitGroup

	// Запуск горутин для расчета квадратов
	for _, num := range numbers {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			square := n * n

			// Блокируем доступ к общей переменной для добавления результата
			mu.Lock()
			sum += square
			mu.Unlock()
		}(num)
	}

	// Ожидаем завершения всех горутин
	wg.Wait()

	// Вывод результата
	fmt.Println("Сумма квадратов:", sum)
}
//////////////////////////////////////
/////////////////////////////////////
// 2 Решение
func main() {
	numbers := []int{2, 4, 6, 8, 10}
	results := make(chan int, len(numbers))

	for _, num := range numbers {
		go func(n int) {
			results <- n * n
		}(num)
	}

	var sum int
	for i := 0; i < len(numbers); i++ {
		sum += <-results
	}

	fmt.Println("Сумма квадратов:", sum)
}
// 4 Реализовать постоянную запись данных в канал (главный поток). Реализовать набор из N воркеров,
//  которые читают произвольные данные из канала и выводят в stdout. Необходима возможность выбора
//  количества воркеров при старте. Программа должна завершаться по нажатию Ctrl+C. Выбрать и обосновать 
// способ завершения работы всех воркеров.

func worker(ctx context.Context, id int, dataChan <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done(): // Завершение по сигналу
			fmt.Printf("Worker %d exiting\n", id)
			return
		case data, ok := <-dataChan: // Чтение из канала
			if !ok {
				return
			}
			fmt.Printf("Worker %d received: %d\n", id, data)
		}
	}
}

func main() {
	// Считываем количество воркеров из аргументов или по умолчанию
	numWorkers := 3
	if len(os.Args) > 1 {
		if n, err := strconv.Atoi(os.Args[1]); err == nil {
			numWorkers = n
		}
	}

	dataChan := make(chan int)
	var wg sync.WaitGroup

	// Контекст для завершения работы
	ctx, cancel := context.WithCancel(context.Background())

	// Запуск воркеров
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, dataChan, &wg)
	}

	// Генерация данных
	go func() {
		i := 0
		for {
			select {
			case <-ctx.Done(): // Завершение генератора
				close(dataChan)
				return
			default:
				dataChan <- i
				i++
				time.Sleep(500 * time.Millisecond) // Задержка
			}
		}
	}()

	// Ожидание сигнала завершения (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	// Завершаем контекст
	cancel()
	wg.Wait()
	fmt.Println("All workers stopped")
}
//////////////////////////////////////////////
/////////////////////////////////////////////
///////////////////////////////////////////////
// 5 Разработать программу, которая будет последовательно отправлять значения в канал, а с другой
//  стороны канала — читать. По истечению N секунд программа должна завершаться.
package main

import (
	"fmt"
	"time"
)

func main() {
	N := 8 * time.Second
	ch := make(chan int)

	// Горутин для записи в канал
	go func() {
		i := 0
		for {
			ch <- i
			i++
			time.Sleep(500 * time.Millisecond) // Имитация задержки
		}
	}()

	// Таймер на N секунд
	timeout := time.After(N)

	for {
		select {
		case val := <-ch:
			fmt.Println("Прочитано:", val)
		case <-timeout:
			fmt.Println("Время вышло. Завершение.")
			return
		}
	}
}
//////////////////////////////////////////////////////////
////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
// 2 Решение
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	N := 8 * time.Second
	ch := make(chan int)

	// Контекст с ручным завершением
	ctx, cancel := context.WithCancel(context.Background())

	// Горутин для записи
	go func() {
		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- i
				i++
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// Таймер для завершения
	go func() {
		time.Sleep(N)
		cancel()
		close(ch)
	}()

	// Чтение из канала
	for val := range ch {
		fmt.Println("Прочитано:", val)
	}

	fmt.Println("Завершение программы.")
}
//////////////////////////////////////////////////
//////////////////////////////////////////////////
//////////////////////////////////////////////////
// 6 Реализовать все возможные способы остановки выполнения горутины.
package main

import (
	"fmt"
	"time"
)

func main() {
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("Горутина завершена через канал")
				return
			default:
				fmt.Println("Работа горутины...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	time.Sleep(2 * time.Second)
	stop <- true // Отправляем сигнал остановки
	time.Sleep(1 * time.Second)
}
///////////////////////////////////////////////
//////////////////////////////////////////////
// 2 Решение 
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Горутина завершена через контекст")
				return
			default:
				fmt.Println("Работа горутины...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}(ctx)

	time.Sleep(2 * time.Second)
	cancel() // Отправляем сигнал отмены
	time.Sleep(1 * time.Second)
}
/////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////
// 3 Решение
package main

import (
	"fmt"
	"time"
)

func main() {
	stop := time.After(2 * time.Second)

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("Горутина завершена через таймер")
				return
			default:
				fmt.Println("Работа горутины...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	time.Sleep(3 * time.Second)
}
///////////////////////////////////////////////
//////////////////////////////////////////////
// 4 Решение
package main

import (
	"fmt"
	"time"
)

func main() {
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("Горутина завершена через закрытие канала")
				return
			default:
				fmt.Println("Работа горутины...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	time.Sleep(2 * time.Second)
	close(stop) // Закрываем канал
	time.Sleep(1 * time.Second)
}
////////////////////////////////////////////////
///////////////////////////////////////////////
// 7 Реализовать конкурентную запись данных в map.
package main

import (
	"fmt"
	"sync"
)

func main() {
	var mu sync.Mutex
	m := make(map[int]int)

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mu.Lock()
			m[i] = i * i
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	fmt.Println("Результат:", m)
}
/////////////////////////////////////////////
////////////////////////////////////////////
/// 2 Решение
package main

import (
	"fmt"
	"sync"
)

func main() {
	var m sync.Map

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Store(i, i*i)
		}(i)
	}
	wg.Wait()

	// Чтение данных
	m.Range(func(key, value any) bool {
		fmt.Printf("Ключ: %v, Значение: %v\n", key, value)
		return true
	})
}
/////////////////////////////////////////////
//////////////////////////////////////////////
// 3 Решение
package main

import (
	"fmt"
)

func main() {
	m := make(map[int]int)
	ch := make(chan func())

	// Горутин для обработки операций с картой
	go func() {
		for f := range ch {
			f()
		}
	}()

	// Конкурентная запись
	for i := 0; i < 10; i++ {
		i := i // Локальная копия переменной
		ch <- func() {
			m[i] = i * i
		}
	}

	// Закрытие канала и завершение
	close(ch)

	// Вывод результата
	fmt.Println("Результат:", m)
}
///////////////////////////////////////////
///////////////////////////////////////////
// 8 Дана переменная int64. Разработать программу которая устанавливает i-й бит в 1 или 0.
package main

import (
	"fmt"
)

func setBits(num int64, indices []uint, set bool) int64 {
	for _, i := range indices {
		if set {
			num |= (1 << i) // Устанавливаем в 1
		} else {
			num &^= (1 << i) // Устанавливаем в 0
		}
	}
	return num
}

func main() {
	var num int64 = 0
	indices := []uint{1, 3, 5} // Устанавливаем 1-й, 3-й и 5-й биты

	num = setBits(num, indices, true) // Установить в 1
	fmt.Printf("Результат установки: %064b\n", num)

	num = setBits(num, indices, false) // Установить в 0
	fmt.Printf("Результат сброса: %064b\n", num)
}
//////////////////////////////////////////////////
/////////////////////////////////////////////////
// 9 Разработать конвейер чисел. Даны два канала:
// в первый пишутся числа (x) из массива, во второй — результат
// операции x*2, после чего данные из второго канала должны выводиться в stdout.
package main

import (
	"fmt"
)

func main() {
	numbers := []int{1, 2, 3, 4, 5}
	in := make(chan int)
	out := make(chan int)

	// Горутин для записи чисел в канал
	go func() {
		for _, num := range numbers {
			in <- num
		}
		close(in)
	}()

	// Горутин для умножения чисел на 2
	go func() {
		for num := range in {
			out <- num * 2
		}
		close(out)
	}()

	// Чтение данных из выходного канала
	for result := range out {
		fmt.Println(result)
	}
}
////////////////////////////////////////
// 2 Решение
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	numbers := []int{1, 2, 3, 4, 5}
	in := make(chan int)
	out := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Генерация чисел
	go func() {
		for _, num := range numbers {
			select {
			case <-ctx.Done():
				close(in)
				return
			case in <- num:
			}
		}
		close(in)
	}()

	// Умножение чисел на 2
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(out)
				return
			case num, ok := <-in:
				if !ok {
					close(out)
					return
				}
				out <- num * 2
			}
		}
	}()

	// Чтение результатов
	for result := range out {
		fmt.Println(result)
	}
}
////////////////////////////////////////////
// 2 Решение
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	numbers := []int{1, 2, 3, 4, 5}
	in := make(chan int)
	out := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Генерация чисел
	go func() {
		for _, num := range numbers {
			select {
			case <-ctx.Done():
				close(in)
				return
			case in <- num:
			}
		}
		close(in)
	}()

	// Умножение чисел на 2
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(out)
				return
			case num, ok := <-in:
				if !ok {
					close(out)
					return
				}
				out <- num * 2
			}
		}
	}()

	// Чтение результатов
	for result := range out {
		fmt.Println(result)
	}
}
////////////////////////////////////////////
////////////////////////////////////////////
// 10 Дана последовательность температурных колебаний:
// -25.4, -27.0 13.0, 19.0, 15.5, 24.5, -21.0, 32.5.
// Объединить данные значения в группы с шагом в 10 градусов.
// Последовательность в подмножноствах не важна.
//Пример: -20:{-25.0, -27.0, -21.0}, 10:{13.0, 19.0, 15.5}, 20: {24.5}, etc.
package main

import (
	"fmt"
)

func main() {
	temperatures := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	groups := make(map[int][]float64)

	for _, temp := range temperatures {
		key := int(temp/10) * 10 // Определяем ключ группы
		groups[key] = append(groups[key], temp)
	}

	// Вывод результата
	for key, temps := range groups {
		fmt.Printf("%d: %v\n", key, temps)
	}
}
/////////////////////////////////////////
////////////////////////////////////////
// 11 Реализовать пересечение двух неупорядоченных множеств.
package main

import "fmt"

func intersection(set1, set2 []int) []int {
	setMap := make(map[int]bool)
	result := []int{}

	// Добавляем элементы первого множества в map
	for _, v := range set1 {
		setMap[v] = true
	}

	// Проверяем пересечение со вторым множеством
	for _, v := range set2 {
		if setMap[v] {
			result = append(result, v)
			delete(setMap, v) // Удаляем, чтобы избежать дублирования
		}
	}

	return result
}

func main() {
	set1 := []int{1, 2, 3, 4, 5}
	set2 := []int{3, 4, 5, 6, 7}
	fmt.Println(intersection(set1, set2)) // [3 4 5]
}
////////////////////////////////////////
///////////////////////////////////////////////////
// 2 Решение 
package main

import (
	"fmt"
	"sync"
)

func intersection(set1, set2 []int) []int {
	setMap := make(map[int]struct{})
	result := []int{}
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Заполняем map параллельно
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, v := range set1 {
			mu.Lock()
			setMap[v] = struct{}{}
			mu.Unlock()
		}
	}()

	// Находим пересечение параллельно
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, v := range set2 {
			mu.Lock()
			if _, exists := setMap[v]; exists {
				result = append(result, v)
				delete(setMap, v)
			}
			mu.Unlock()
		}
	}()

	wg.Wait()
	return result
}

func main() {
	set1 := []int{1, 2, 3, 4, 5}
	set2 := []int{3, 4, 5, 6, 7}
	fmt.Println(intersection(set1, set2)) // [3 4 5]
}
/////////////////////////////////////////////////
// 3 Решение 
package main

import (
	"fmt"
	"sort"
)

func intersection(set1, set2 []int) []int {
	sort.Ints(set1)
	sort.Ints(set2)

	i, j := 0, 0
	result := []int{}

	for i < len(set1) && j < len(set2) {
		if set1[i] == set2[j] {
			result = append(result, set1[i])
			i++
			j++
		} else if set1[i] < set2[j] {
			i++
		} else {
			j++
		}
	}

	return result
}

func main() {
	set1 := []int{1, 2, 3, 4, 5}
	set2 := []int{3, 4, 5, 6, 7}
	fmt.Println(intersection(set1, set2)) // [3 4 5]
}
/////////////////////////////////////////////////////
// 12 Имеется последовательность строк - (cat, cat, dog, cat, tree) создать для нее собственное множество.
package main

import "fmt"

func createSet(input []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, v := range input {
		set[v] = struct{}{}
	}
	return set
}

func main() {
	data := []string{"cat", "cat", "dog", "cat", "tree"}
	set := createSet(data)
	fmt.Println(set) // map[cat:{} dog:{} tree:{}]
}
///////////////////////////////////////////////////
// 13 Поменять местами два числа без создания временной переменной.
package main

import "fmt"

func main() {
	a, b := 5, 10
	a = a + b
	b = a - b
	a = a - b
	fmt.Println(a, b) // 10, 5
}
//////////////////////////////////////////////////
package main

import "fmt"

func main() {
	a, b := 5, 10
	a = a * b
	b = a / b
	a = a / b
	fmt.Println(a, b) // 10, 5
}
/////////////////////////////////////////////
package main

import "fmt"

func main() {
	a, b := 5, 10
	a, b = b, a
	fmt.Println(a, b) // 10, 5
}
// 14 Разработать программу, которая в рантайме способна определить
//  тип переменной: int, string, bool, channel из переменной типа interface{}.
package main

import "fmt"

func determineType(v interface{}) {
	switch v.(type) {
	case int:
		fmt.Println("Type is int")
	case string:
		fmt.Println("Type is string")
	case bool:
		fmt.Println("Type is bool")
	case chan int:
		fmt.Println("Type is channel of int")
	default:
		fmt.Println("Unknown type")
	}
}

func main() {
	var a interface{}
	determineType(a)
}
// 15 К каким негативным последствиям может привести данный фрагмент кода, 
// и как это исправить? Приведите корректный пример реализации.


// var justString string
// func someFunc() {
//   v := createHugeString(1 << 10)
//   justString = v[:100]
// }

// func main() {
//   someFunc()
// }
package main

import "fmt"

var justString string

func createHugeString(size int) string {
	return string(make([]byte, size))
}

func someFunc() {
	v := createHugeString(1 << 10)
	// Копируем первые 100 символов, чтобы предотвратить утечку памяти
	justString = string([]byte(v[:100]))
}

func main() {
	someFunc()
	fmt.Println(len(justString)) // 100
}

// 16 Реализовать быструю сортировку массива (quicksort) встроенными методами языка.
package main

import (
	"fmt"
	"sort"
)

func main() {
	arr := []int{3, 6, 8, 10, 1, 2, 1}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})
	fmt.Println(arr) // [1 1 2 3 6 8 10]
}
//////////////////////////////////////////////////
/////////////////////////////////////////////////
package main

import (
	"fmt"
	"sort"
)

func main() {
	arr := []int{3, 6, 8, 10, 1, 2, 1}
	sort.Ints(arr)
	fmt.Println(arr) // [1 1 2 3 6 8 10]
}
////////////////////////////////////////
///////////////////////////////////////
package main

import "fmt"

func quickSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}

	pivot := arr[0]
	var less, greater []int
	for _, val := range arr[1:] {
		if val <= pivot {
			less = append(less, val)
		} else {
			greater = append(greater, val)
		}
	}

	return append(append(quickSort(less), pivot), quickSort(greater)...)
}

func main() {
	var arr []int
	fmt.Println(quickSort(arr)) // [1 1 2 3 6 8 10]
}
///////////////////////////////////////////
// 17 Реализовать бинарный поиск встроенными методами языка.
package main

import "fmt"

func binarySearch(arr []int, x int) int {
	low, high := 0, len(arr)-1

	for low <= high {
		mid := (low + high) / 2
		if arr[mid] == x {
			return mid
		}
		if arr[mid] < x {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1 // элемент не найден
}

func main() {

}
//////////////////////////////////////////////////
package main

import (
	"fmt"
	"sort"
)

func main() {
	arr := []int{1, 3, 5, 7, 9}
	x := 5
	index := sort.Search(len(arr), func(i int) bool {
		return arr[i] >= x
	})

	if index < len(arr) && arr[index] == x {
		fmt.Printf("Found %d at index %d\n", x, index) // Found 5 at index 2
	} else {
		fmt.Printf("%d not found\n", x)
	}
}
////////////////////////////////////////////////////
//18 Реализовать структуру-счетчик, которая будет инкрементироваться в 
// конкурентной среде. По завершению программа должна выводить итоговое значение счетчика.
package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	var wg sync.WaitGroup
	counter := &Counter{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Println("Final Counter Value:", counter.Value()) // Final Counter Value: 1000
}
////////////////////////////////////////////////////////
/////////////////////////////////////////////////////
package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	ch    chan struct{}
	value int
}

func NewCounter() *Counter {
	return &Counter{ch: make(chan struct{}, 1)}
}

func (c *Counter) Increment() {
	c.ch <- struct{}{} // Заблокировать
	c.value++
	<-c.ch // Разблокировать
}

func (c *Counter) Value() int {
	c.ch <- struct{}{}
	val := c.value
	<-c.ch
	return val
}

func main() {
	var wg sync.WaitGroup
	counter := NewCounter()

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Println("Final Counter Value:", counter.Value()) // Final Counter Value: 1000
}
/////////////////////////////////////////////////////////
//19 Разработать программу, которая переворачивает подаваемую 
// на ход строку (например: «главрыба — абырвалг»). Символы могут быть unicode.
package main

import (
	"fmt"
)

func reverseString(input string) string {
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func main() {

}
//////////////////////////////////////////////////////////////
package main

import (
	"fmt"
)

func reverseString(input string) string {
	runes := []rune(input)
	result := make([]rune, len(runes))
	for i := len(runes) - 1; i >= 0; i-- {
		result[len(runes)-1-i] = runes[i]
	}
	return string(result)
}
/////////////////////////////////////////////////
func reverseString(input string) string {
	runes := []rune(input)
	if len(runes) == 0 {
		return ""
	}
	return string(runes[len(runes)-1]) + reverseString(string(runes[:len(runes)-1]))
}
////////////////////////////////////////////
// Разработать программу, которая переворачивает слова в строке. 
// Пример: «snow dog sun — sun dog snow».
package main

import (
	"fmt"
	"strings"
)

func reverseWords(input string) string {
	words := strings.Fields(input)
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}
	return strings.Join(words, " ")
}
// //////////////////////////////////
func reverseWords(input string) string {
	words := strings.Fields(input)
	stack := make([]string, len(words))
	copy(stack, words)

	var result []string
	for len(stack) > 0 {
		n := len(stack)
		result = append(result, stack[n-1])
		stack = stack[:n-1]
	}
	return strings.Join(result, " ")
}
// ////////////////////////////////////////
//  21 Реализовать паттерн «адаптер» на любом примере.
package main

import "fmt"

// Старый интерфейс
type OldPrinter interface {
	PrintOld(message string)
}

// Реализация старого интерфейса
type OldPrinterImpl struct{}

func (o *OldPrinterImpl) PrintOld(message string) {
	fmt.Println("Old Printer: " + message)
}

// Новый интерфейс
type NewPrinter interface {
	PrintNew(message string)
}

// Адаптер, который работает с новым интерфейсом
type PrinterAdapter struct {
	oldPrinter OldPrinter
}

func (p *PrinterAdapter) PrintNew(message string) {
	p.oldPrinter.PrintOld("Adapted: " + message)
}

func main() {

}
/////////////////////////////////////////////////////
// 22 Разработать программу, которая перемножает, делит, складывает,
//  вычитает две числовых переменных a,b, значение которых > 2^20.
package main

import (
	"fmt"
	"math/big"
)

func main() {
	// Создаем большие числа
	a := big.NewInt(1 << 21) // a = 2^21
	b := big.NewInt(1 << 22) // b = 2^22

	// Сложение
	sum := new(big.Int).Add(a, b)

	// Вычитание
	sub := new(big.Int).Sub(a, b)


	// Умножение
	mul := new(big.Int).Mul(a, b)


	// Деление
	div := new(big.Int).Div(b, a)

}
///////////////////////////////////////
///////////////////////////////////
// 23 Удалить i-ый элемент из слайса.
package main

import "fmt"

func main() {
	slice := []int{1, 2, 3, 4, 5}
	i := 2 

	slice = append(slice[:i], slice[i+1:]...)

}
////////////////////////////////////////
func main() {
	slice := []int{1, 2, 3, 4, 5}
	i := 2 // Удаляем элемент с индексом 2

	copy(slice[i:], slice[i+1:])
	slice = slice[:len(slice)-1]
	fmt.Println(slice) // Output: [1 2 4 5]
}
//////////////////////////////////////////
// 24 Разработать программу нахождения расстояния между
//  двумя точками, которые представлены в виде структуры Point
//  с инкапсулированными параметрами x,y и конструктором.
package main

import (
	"fmt"
	"math"
)

// Point структура с инкапсулированными координатами
type Point struct {
	x, y float64
}

// Конструктор
func NewPoint(x, y float64) Point {
	return Point{x: x, y: y}
}

// Метод для вычисления расстояния до другой точки
func (p Point) DistanceTo(other Point) float64 {
	return math.Sqrt(math.Pow(p.x-other.x, 2) + math.Pow(p.y-other.y, 2))
}

func main() {
	p1 := NewPoint(1.0, 2.0)
	p2 := NewPoint(4.0, 6.0)

	fmt.Println("Distance:", p1.DistanceTo(p2)) // Output: 5
}
///////////////////////////////////////////////////////
//  25 Реализовать собственную функцию sleep.
package main

import (
	"fmt"
	"time"
)

func Sleep(duration time.Duration) {
	<-time.After(duration)
}

func main() {
	fmt.Println("Start")
	Sleep(2 * time.Second)
	fmt.Println("End")
}
/////////////////////////////////////////////////
///////////////////////////////////////////
package main

import (
	"fmt"
	"time"
)

func Sleep(duration time.Duration) {
	start := time.Now()
	for time.Since(start) < duration {
		// Пустой цикл
	}
}

func main() {
	fmt.Println("Start")
	Sleep(2 * time.Second)
	fmt.Println("End")
}
////////////////////////////////////////////////////
// Разработать программу, которая проверяет, что все символы в строке
//  уникальные (true — если уникальные, false etc). Функция проверки должна быть регистронезависимой.

// Например: 
// abcd — true
// abCdefAaf — false
	// aabcd — false
	package main

	import (
		"fmt"
		"strings"
	)
	
	func IsUnique(s string) bool {
		charMap := make(map[rune]bool)
		s = strings.ToLower(s)
		for _, char := range s {
			if charMap[char] {
				return false
			}
			charMap[char] = true
		}
		return true
	}
	
	func main() {

	}
///////////////////////////////////////////
package main

import (
	"fmt"
	"strings"
)

func IsUnique(s string) bool {
	unique := ""
	s = strings.ToLower(s)
	for _, char := range s {
		if strings.ContainsRune(unique, char) {
			return false
		}
		unique += string(char)
	}
	return true
}

func main() {

}
