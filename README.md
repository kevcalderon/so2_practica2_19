**<h1 align="center">PRACTICA 3</h1>**

**<h3 align="center">Control y Monitoreo de Procesos y Memoria Fase 3</h3>**

### Datos de los estudiantes

<div align="center">

| Nombre                      | Carné     |
| --------------------------- | --------- |
| Adrian Samuel Molina Cabrera    | 201903850 |
| Carlos Raúl Campos Meléndez | 201800639 |
| Kevin Josué Calderón Peraza | 201902714 |

</div>


# MANUAL TECNICO
* Descripción de la configuración de la Máquina virtual
  Estas son las configuraciones que se realizaron para la creacion de maquina virtual asi como el usuario.
  
<p align="center">
  <img src="https://github.com/kevcalderon/so2_practica2_19/blob/main/imgs/direccion.jpg" width="600">
</p>
<p align="center">
  <img src="https://github.com/kevcalderon/so2_practica2_19/blob/main/imgs/conf.jpg" width="600">
</p>
<p align="center">
  <img src="https://github.com/kevcalderon/so2_practica2_19/blob/main/imgs/grub.jpg" width="600">
</p>
<p align="center">
  <img src="https://github.com/kevcalderon/so2_practica2_19/blob/main/imgs/users.jpg" width="600">
</p>

* Descripción de la configuración de la Máquina virtual
<p align="center">
  <img src="https://github.com/kevcalderon/SO1_201902714/blob/master/Practica2/img/vms.png" width="600">
</p>

## /proc/[pid]/maps
Este archivo muestra los rangos de memoria virtual asignados a un proceso específico, junto con los permisos de acceso y la ubicación de los archivos mapeados en esos rangos. Cada línea del archivo representa un rango de memoria y contiene información como dirección inicial, dirección final, permisos y desplazamiento en disco.

## /proc/[pid]/smaps
Este archivo proporciona información más detallada sobre el uso de memoria de un proceso en particular. Contiene estadísticas sobre diferentes regiones de memoria, como tamaño, dirección, permisos, uso compartido, bibliotecas cargadas, información de archivo asociada y consumo de memoria física.

Ambos archivos son útiles para realizar análisis y depuración de problemas relacionados con la memoria en Linux. Pueden ayudarte a identificar fugas de memoria, analizar el uso de la memoria por parte de un proceso específico y comprender cómo se está asignando y utilizando la memoria en general en el sistema.

# Comandos de configuracion para ArchLinux

    fdisk -l

    ls /sys/firmware/efi/efivars

    fdisk /dev/sda

## Crear nueva particion
    n                               
    default
    default
    2048
    4G

## Crear nueva particion
    n                               
    default
    default
    default
    default

## Guardar la configuracion
    w                               


## formatear la particion
    mkswap /dev/sda1                
    mkfs.ext4 /dev/sda2

    mount /dev/sda2 /mnt
    swapon /dev/sda1

## descargar los paquetes
    pacstrap /mnt base base-devel

    genfstab -U /mnt >> /mnt/etc/fstab
    arch-chroot /mnt

## configurar reloj
    ln -sf /usr/share/zoneinfo/Mexico/BajaSur /etc/localtime

## instalar nano
    pacman -S nano                  

## agregar al final del archivo las lineas
    nano /etc/pacman.conf           
        [kernel-lts]
        Server = https://repo.m2x.dev/current/$repo/$arch

## se agrega la llave
    pacman-key --recv-keys 76C6E477042BFE985CC2209C08A255442FAFF0  
    pacman-key --finger 76C6E477042BFE985CC2209C08A255442FAFF0
    pacman-key --lsign-key 76C6E477042BFE985CC2209C08A255442FAFF0

## Se instalan los headers
    pacman -Syu linux-lts54 linux-lts54-headers     
        default


## configuracion del grub
    pacman -S grub
    grub-install /dev/sda
    grub-mkconfig -o /boot/grub/grub.cfg


## se editan los lenguajes que se necesiten
    nano /etc/locale.get
    locale-gen
    echo LANG=en_US.UTF-8 > /etc/locale.conf
    export LANG=en_US.UTF-8
    echo myarch > /etc/hostname

## Se configura la red
    touch /etc/hosts    
    nano /etc/hosts

        127.0.0.1   localhost
        ::1         localhost
        127.0.1.1   myarch

## Se configuran los usuarios
    passwd
    pacman -S sudo
    useradd -m so2_Ejemplo
    passwd so2_Ejemplo 
    usermod -aG wheel,audio,video,storage so2_Ejemplo


    pacman -S xorg networkmanager
        default

## instalacion de interfaz grafica
    pacman -S gnome     
        default
        default

    systemctl enable gdm.service
    umount /mnt

## instalacion de Golang
    sudo pacman -S wget
    wget https://go.dev/dl/go1.20.4.linux-amd64.tar.gz
    sudo tar -xvf go1.20.4.linux-amd64.tar.gz -C /usr/local/
    export PATH=$PATH:/usr/local/go/bin

# MANUAL DE USUARIO
* El proyecto fue realizado con la tecnologia de react, go, modulos de kernel con la finalidad de ser mas dinamico las funcionalidades.
* En este se puede visualizar un poligono de frecuencia en el cual indica la ram consumida, asi como tambien los procesos del kernel

  
<p align="center">
  <img src="https://github.com/kevcalderon/so2_practica2_19/blob/main/imgs/front1.jpg" width="600">
</p>


* Cada proceso se puede tiene una opcion de visualizar los nodos hijos por medio de un boton y por medio de un popup mostralos.
* Asi como tambien un botón de kill para matar el proceso.
* Tambien, un boton para visualizar de manera especifica un proceso.

![cd5c245f-b6e1-413f-a57c-b4defd59c039](https://github.com/kevcalderon/so2_practica2_19/assets/18681538/a5fb34e0-fc95-4a9e-a03f-cfc745f8ac68)

<p align="center">
  <img src="https://github.com/kevcalderon/so2_practica2_19/blob/main/imgs/front2.jpg" width="600">
</p>


* En esta vista se muestran mas especificos del proceso.
*  Como por ejemplo RSS: Esta línea muestra la cantidad de memoria física utilizada actualmente por el proceso. El valor se muestra en Megabytes (MB). Y Tambien Size: Esta línea muestra el tamaño total de la memoria virtual reservada por el proceso. El valor también se muestra en Megabytes (MB).
*  Tambien se puede visualizar un mapa de la memoria, este incluye el inicio del segmento el fin del segmento, memoria residente y memoria virtual del proceso.

![c09615af-0692-4ac2-9773-c7ad588afa8c](https://github.com/kevcalderon/so2_practica2_19/assets/18681538/594cb297-754e-4f31-a4af-c1fa17c3acf2)
