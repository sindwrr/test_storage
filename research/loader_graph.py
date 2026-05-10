import matplotlib.pyplot as plt

x = ["1 KB", "10 KB", "100 KB", "512 KB", "1 MB", "5 MB", "10 MB", "20 MB", "30 MB"]

y1 = [2.6, 2.7, 3.0, 4.6, 5.0, 11.2, 23.3, 41.8, 51.8]
y2 = [3.7, 3.9, 4.9, 8.0, 12.2, 54.4, 107.9, 217.4, 310.3]
y3 = [7.5, 8.6, 9.8, 18.6, 27.2, 118.2, 203.7, 470.7, 604.3]
y4 = [15.5, 19.9, 21.8, 32.3, 46.8, 218.9, 439.8, 964.8, 1335]
y5 = [27.6, 29.8, 34.4, 44.9, 81.5, 380.2, 635.1, 1120, 1930]

plt.figure(figsize=(10, 6))

plt.plot(x, y1, marker="o", linewidth=2, label="concurrency = 1")
plt.plot(x, y2, marker="o", linewidth=2, label="concurrency = 5")
plt.plot(x, y3, marker="o", linewidth=2, label="concurrency = 10")
plt.plot(x, y4, marker="o", linewidth=2, label="concurrency = 20")
plt.plot(x, y5, marker="o", linewidth=2, label="concurrency = 30")

plt.xlabel("Размер файла", fontsize=12)
plt.ylabel("Средняя задержка, мс", fontsize=12)
plt.title("Зависимость времени отклика от параллельной загрузки в систему", fontsize=14)

plt.grid(True, linestyle="--")
plt.legend(title="Кол-во одновременных запросов", fontsize=10)

plt.axhline(y=2000, color="red", linestyle=":", linewidth=2, label="Лимит 2 с")

plt.tight_layout()
plt.savefig("loader_graph.png")
