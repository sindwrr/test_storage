import matplotlib.pyplot as plt

x = list(i for i in range(1, 31))

y = [
    0.25,
    0.27,
    0.37,
    0.45,
    0.51,
    0.58,
    0.72,
    0.75,
    0.77,
    0.82,
    0.93,
    0.92,
    1.06,
    1.11,
    1.12,
    1.23,
    1.27,
    1.30,
    1.46,
    1.48,
    1.51,
    1.57,
    1.65,
    1.66,
    1.68,
    1.79,
    1.78,
    1.85,
    1.91,
    1.94,
]

plt.figure(figsize=(10, 6))

plt.plot(x, y, marker="o", linewidth=2)

plt.xlabel("Кол-во одновременных запросов", fontsize=12)
plt.ylabel("Потребляемая память, ГБ", fontsize=12)
plt.title("Зависимость потребления RAM от параллельной загрузки в систему", fontsize=14)

plt.grid(True, linestyle="--")

plt.axhline(y=2, color="red", linestyle=":", linewidth=2, label="Лимит 2 ГБ")

plt.tight_layout()
plt.savefig("ram_graph.png")
