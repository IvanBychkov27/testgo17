// https://goprog.ru/posts/go-tensorflow

//TensorFlow написан на Си. Подключение к другим языкам осуществляется через биндинги. Поэтому для работы с TensorFlow в систему должна быть установлена Си-библиотека.
//Существует 64-битная версия версии для Linux (CPU и GPU) и для MacOS.
//Для установки в систему нужно будет распаковать архив в /usr/local/lib или в другую директории (надо не забыть добавить ее в LD_LIBRARY_PATH).
//sudo tar -xz libtensorflow.tar.gz -C /usr/local
//Устанавливаем биндинг для Golang со всеми зависимостями:
//go get github.com/tensorflow/tensorflow/tensorflow/go
//Проверяем работу инсталляции:
//go test github.com/tensorflow/tensorflow/tensorflow/go

package main

//
//import (
//	"fmt"
//	tf "github.com/tensorflow/tensorflow/tensorflow/go"
//
//)
//
//func main() {
//
//	// replace myModel and myTag with the appropriate exported names in the chestrays-keras-binary-classification.ipynb
//	model, err := tf.LoadSavedModel("myModel", []string{"myTag"}, nil)
//
//	if err != nil {
//		fmt.Printf("Error loading saved model: %s\n", err.Error())
//		return
//	}
//
//	defer model.Session.Close()
//
//	tensor, _ := tf.NewTensor([1][250][250][3]float32{})
//
//	result, err := model.Session.Run(
//		map[tf.Output]*tf.Tensor{
//			model.Graph.Operation("inputLayer_input").Output(0): tensor, // Replace this with your input layer name
//		},
//		[]tf.Output{
//			model.Graph.Operation("inferenceLayer/Sigmoid").Output(0), // Replace this with your output layer name
//		},
//		nil,
//	)
//
//	if err != nil {
//		fmt.Printf("Error running the session with input, err: %s\n", err.Error())
//		return
//	}
//
//	fmt.Printf("Result value: %v \n", result[0].Value())
//
//
//
//	fmt.Println("Done...")
//}

//
//import (
//	"fmt"
//
//	tf "github.com/tensorflow/tensorflow/tensorflow/go"
//	"github.com/tensorflow/tensorflow/tensorflow/go/op"
//)
//
//func main() {
//
//	// создаем граф с одной строковой константой
//	s := op.NewScope()
//	c := op.Const(s, "Hello from TensorFlow version " + tf.Version())
//	graph, err := s.Finalize()
//	if err != nil {
//		panic(err)
//	}
//
//	// выполняем вычисления в сессии
//	sess, err := tf.NewSession(graph, nil)
//	if err != nil {
//		panic(err)
//	}
//
//	output, err := sess.Run(nil, []tf.Output{c}, nil)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(output[0].Value())
//}
