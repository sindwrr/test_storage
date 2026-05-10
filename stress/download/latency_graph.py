import matplotlib.pyplot as plt

x = ["1 KB", "10 KB", "100 KB", "512 KB", "1 MB", "5 MB", "10 MB", "20 MB", "30 MB"]

y1 = [0.202, 0.221, 0.289, 0.548, 0.786, 4.5, 11.4, 21.1, 31.2]
y2 = [0.881, 0.891, 0.991, 1.6, 2.2, 12.8, 40.4, 85.1, 131.7]
y3 = [2.1, 2.2, 2.6, 3.5, 4.3, 25.6, 79.8, 174.0, 262.0]
y4 = [5.4, 5.8, 6.8, 8.1, 9.1, 49.7, 156.9, 348.4, 524.0]
y5 = [8.5, 8.7, 9.2, 12.4, 15.4, 77.5, 237.6, 523.9, 788.7]

plt.figure(figsize=(10, 6))

plt.plot(x, y1, marker="o", linewidth=2, label="concurrency = 1")
plt.plot(x, y2, marker="o", linewidth=2, label="concurrency = 5")
plt.plot(x, y3, marker="o", linewidth=2, label="concurrency = 10")
plt.plot(x, y4, marker="o", linewidth=2, label="concurrency = 20")
plt.plot(x, y5, marker="o", linewidth=2, label="concurrency = 30")

plt.xlabel("Размер файла", fontsize=12)
plt.ylabel("Средняя задержка, мс", fontsize=12)
plt.title(
    "Зависимость времени отклика от параллельного скачивания из системы", fontsize=14
)

plt.yscale("log")

plt.grid(True, linestyle="--")
plt.legend(title="Кол-во одновременных запросов", fontsize=10)

plt.axhline(y=2000, color="red", linestyle=":", linewidth=2, label="Лимит 2 с")

plt.tight_layout()
plt.savefig("latency_graph.png")
