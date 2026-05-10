import matplotlib.pyplot as plt

x = list(i for i in range(1, 31))

y = [
    98.4,
    99.3,
    98.16,
    98.87,
    99.68,
    100.2,
    99.42,
    100.3,
    100.4,
    100.8,
    100.3,
    100.7,
    101.7,
    101.4,
    101,
    101.4,
    101.9,
    102.4,
    101.8,
    102.8,
    103.3,
    102.6,
    102.8,
    103.7,
    104.1,
    104.6,
    103.8,
    104.8,
    104.6,
    104.9,
]

plt.figure(figsize=(10, 6))

plt.plot(x, y, marker="o", linewidth=2)

plt.xlabel("Кол-во одновременных запросов", fontsize=12)
plt.ylabel("Потребляемая память, МБ", fontsize=12)
plt.title(
    "Зависимость потребления RAM от параллельного скачивания из системы", fontsize=14
)

plt.grid(True, linestyle="--")

plt.tight_layout()
plt.savefig("ram_graph.png")
