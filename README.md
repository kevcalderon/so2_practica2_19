**<h1 align="center">PRACTICA 2</h1>**

**<h3 align="center">Control y monitoreo de procesos, fase 2</h3>**

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
* El proyecto fue realizado con la tecnologia de react, con la finalidad de ser mas dinamico las funcionalidades.
* En este se puede visualizar un poligono de frecuencia en el cual indica la ram consumida, asi como tambien los procesos del kernel
* Cada proceso tiene una funcionalidad que es ver nodos hijos, eliminarlo y ver especificamente en otra pantalla el consumo de ese proceso.
  
<p align="center">
  <img src="https://github.com/kevcalderon/so2_practica2_19/blob/main/imgs/front1.jpg" width="600">
</p>
<p align="center">
  <img src="https://github.com/kevcalderon/so2_practica2_19/blob/main/imgs/front2.jpg" width="600">
</p>
<p align="center">
  <img src="https://github.com/kevcalderon/so2_practica2_19/blob/main/imgs/front3.jpg" width="600">
</p>
